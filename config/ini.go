package config

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var (
	DEFAULT_SECTION  = "default"   // default section means if some ini items not in a section, make them in default section,
	bNotesFlag       = []byte{'#'} // notes signal
	bSemFlag         = []byte{';'} // semicolon signal
	bEmptyFlag       = []byte{}
	bEqualFlag       = []byte{'='} // equal signal
	bDQuoteFlag      = []byte{'"'} // quote signal
	sectionStartFlag = []byte{'['} // section start signal
	sectionEndFlag   = []byte{']'} // section end signal
	lineBreakFlag    = "\n"
)

type IniConfig struct {
}

type IniConfigManager struct {
	filename       string
	data           map[string]map[string]string // section=> key:val
	sectionComment map[string]string            // section : comment
	keyComment     map[string]string            // id: []{comment, key...}; id 1 is for main comment.
	changed        bool
	sync.RWMutex
}

func (ini *IniConfig) Parse(filename string) (ConfigManager, error) {
	return ini.parseFile(filename)
}

func (ini *IniConfig) parseFile(name string) (ConfigManager, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	cfg := &IniConfigManager{filename: file.Name(),
		data:           make(map[string]map[string]string),
		sectionComment: make(map[string]string),
		keyComment:     make(map[string]string),
		RWMutex:        sync.RWMutex{},
	}
	cfg.Lock()
	defer cfg.Unlock()
	defer file.Close()
	var comment bytes.Buffer //需要将注解记录下来，回写的时候有用
	buf := bufio.NewReader(file)
	head, err := buf.Peek(3)
	if err == nil && head[0] == 239 && head[1] == 187 && head[2] == 191 {
		for i := 1; i <= 3; i++ {
			buf.ReadByte()
		}
	}
	section := DEFAULT_SECTION
	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if bytes.Equal(line, bEmptyFlag) {
			continue
		}
		line = bytes.TrimSpace(line)

		var bComment bool = false
		switch {
		case bytes.HasPrefix(line, bNotesFlag):
			bComment = true

		}
		if bComment {
			line = bytes.TrimLeft(line, string(bNotesFlag))
			line = bytes.TrimLeftFunc(line, unicode.IsSpace)
			comment.Write(line)
			comment.WriteByte('\n')
			continue
		}

		if bytes.HasPrefix(line, sectionStartFlag) && bytes.HasSuffix(line, sectionEndFlag) {
			//new section
			section = string(line[1 : len(line)-1])
			if comment.Len() > 0 {
				cfg.sectionComment[section] = comment.String()
				comment.Reset()
			}
			if _, ok := cfg.data[section]; !ok {
				cfg.data[section] = make(map[string]string)
			}
			continue
		}

		if _, ok := cfg.data[section]; !ok {
			cfg.data[section] = make(map[string]string)
		}

		keyValue := bytes.SplitN(line, bEqualFlag, 2)

		key := string(bytes.TrimSpace(keyValue[0])) // key name case insensitive
		//key = strings.ToLower(key)
		if len(keyValue) != 2 {
			return nil, errors.New("read the content error: \"" + string(line) + "\", should key = val")
		}
		val := bytes.TrimSpace(keyValue[1])
		if bytes.HasPrefix(val, bDQuoteFlag) {
			val = bytes.Trim(val, `"`)
		}
		cfg.data[section][key] = string(val)
		if comment.Len() > 0 {
			cfg.keyComment[section+"::"+key] = comment.String()
			comment.Reset()
		}

	}

	return cfg, nil

}

func (c *IniConfigManager) Set(key, val string) error {
	c.Lock()
	defer c.Unlock()
	if len(key) == 0 {
		return errors.New("key is empty")
	}

	var section, k string
	var sectionKey []string = strings.Split(key, "::")

	if len(sectionKey) >= 2 {
		section = sectionKey[0]
		k = sectionKey[1]
	} else {
		section = DEFAULT_SECTION
		k = sectionKey[0]
	}

	if _, ok := c.data[section]; !ok {
		c.data[section] = make(map[string]string)
	}
	c.data[section][k] = val
	c.changed = true
	return nil
}

func (c *IniConfigManager) String(key string) string {
	return c.getdata(key)
}

func (c *IniConfigManager) Strings(key string) []string {
	return strings.Split(c.getdata(key), ";")
}

func (c *IniConfigManager) Int(key string) (int, error) {
	return strconv.Atoi(c.getdata(key))
}

func (c *IniConfigManager) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.getdata(key), 10, 64)
}

func (c *IniConfigManager) Bool(key string) (bool, error) {
	return strconv.ParseBool(c.getdata(key))
}

func (c *IniConfigManager) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.getdata(key), 64)
}

func (c *IniConfigManager) SaveConfigFile(filename string, bsort bool) error {
	if !c.changed {
		return nil
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := bytes.NewBuffer(nil)

	keys := make([]string, 0, 32)

	for key, _ := range c.data {
		keys = append(keys, key)
	}

	if bsort {
		sort.Strings(keys)
	}

	for _, secKey := range keys {
		section := secKey
		dt := c.data[section]
		//section, dt := range c.data
		// Write section comments.
		if v, ok := c.sectionComment[section]; ok {
			if _, err = buf.WriteString(string(bNotesFlag) + v + lineBreakFlag); err != nil {
				return err
			}
		}

		if section != DEFAULT_SECTION {
			// Write section name.
			if _, err = buf.WriteString(string(sectionStartFlag) + section + string(sectionEndFlag) + lineBreakFlag); err != nil {
				return err
			}
		}

		tmpKeys := make([]string, 0, len(dt))
		for k, _ := range dt {
			tmpKeys = append(tmpKeys, k)
		}

		if bsort {
			sort.Strings(tmpKeys)
		}

		for _, key := range tmpKeys {
			val := dt[key]
			if key != " " {
				// Write key comments.
				if v, ok := c.keyComment[key]; ok {
					if _, err = buf.WriteString(string(bNotesFlag) + v + lineBreakFlag); err != nil {
						return err
					}
				}

				// Write key and value.
				if _, err = buf.WriteString(key + string(bEqualFlag) + val + lineBreakFlag); err != nil {
					return err
				}
			}
		}

		// Put a line between sections.
		if _, err = buf.WriteString(lineBreakFlag); err != nil {
			return err
		}
	}

	if _, err = buf.WriteTo(f); err != nil {
		return err
	}
	c.changed = false
	return nil

}

// section::key or key
func (c *IniConfigManager) getdata(key string) string {
	if len(key) == 0 {
		return ""
	}
	c.RLock()
	defer c.RUnlock()

	var (
		section, k string
		sectionKey []string = strings.Split(key, "::")
	)
	if len(sectionKey) >= 2 {
		section = sectionKey[0]
		k = sectionKey[1]
	} else {
		section = DEFAULT_SECTION
		k = sectionKey[0]
	}
	if v, ok := c.data[section]; ok {
		if vv, ok2 := v[k]; ok2 {
			return vv
		}
	}
	return ""
}

func init() {
	Register("ini", &IniConfig{})
}
