package diff

import (
	"fmt"
	"io"
	"os"

	"github.com/r3labs/diff/v3"
	"gopkg.in/yaml.v2"
)

var result map[string]interface{}

// CheckErr 检查错误并处理
func CheckErr(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// Run 比较两个 YAML 文件并打印差异
func Run(files []string) {
	// 检查文件数量
	if len(files) != 2 {
		CheckErr(fmt.Errorf("需要提供两个文件名进行比较"))
	}

	// 解析文件内容
	yaml1, err := unmarshal(files[0])
	CheckErr(err)
	yaml2, err := unmarshal(files[1])
	CheckErr(err)

	// TODO: 如果出错了，那么输出yaml2中的值
	result = yaml2

	// 计算差异并输出
	result = computeDiff(yaml1, yaml2)
	outputYAML(result)
}

// 解析YAML文件
func unmarshal(filename string) (map[string]interface{}, error) {
	var err error
	var contents []byte

	if filename == "-" {
		contents, err = io.ReadAll(os.Stdin)
	} else {
		contents, err = os.ReadFile(filename)
	}
	if err != nil {
		return nil, fmt.Errorf("无法读取文件 %s: %v", filename, err)
	}
	data := make(map[string]interface{})
	if err := yaml.Unmarshal(contents, &data); err != nil {
		return nil, fmt.Errorf("无法解析文件 %s: %v", filename, err)
	}
	return data, nil
}

// 计算两个YAML结构之间的差异
func computeDiff(a, b map[string]interface{}) map[string]interface{} {
	// 初始化差异映射
	out := map[string]interface{}{}

	// 遍历第二个YAML结构的键值对
	for bk, bv := range b {
		// 检查键值对是否为子映射
		if subMapB, ok := bv.(map[interface{}]interface{}); ok {
			// 检查第一个YAML结构是否存在当前键
			if av, ok := a[bk]; ok {
				// 如果第一个YAML结构中存在与当前键相同的键，且值都是子映射，则递归计算这两个子映射的差异
				if subMapA, ok := av.(map[interface{}]interface{}); ok {
					subDiff := computeDiff(convertMapStringInterface(subMapA), convertMapStringInterface(subMapB))
					// 如果子映射的差异不为空，说明这两个子映射有差异，我们将这个差异加入输出映射中
					if subDiff != nil && len(subDiff) != 0 {
						out[bk] = subDiff
						continue
					}
				}
			}
		}
		// 如果当前键值对不是子映射，则比较第一个YAML结构中当前键对应的值与第二个YAML结构中当前键对应的值之间的差异
		if a[bk] == nil && bv == nil {
			continue
		}
		// 比较差异，允许类型不匹配
		changelog, err := diff.Diff(a[bk], bv, diff.AllowTypeMismatch(true))
		CheckErr(err)
		// 处理计算得到的差异
		if len(changelog) != 0 {
			out[bk] = bv
		}
	}

	// 返回差异映射
	return out
}

// 输出YAML格式的数据
func outputYAML(data map[string]interface{}) {
	yamlData, err := yaml.Marshal(data)
	CheckErr(err)
	fmt.Println(string(yamlData))
}

// 将map[interface{}]interface{}转换为map[string]interface{}
func convertMapStringInterface(m map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[fmt.Sprintf("%v", k)] = v
	}
	return result
}
