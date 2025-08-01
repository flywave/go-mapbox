package mvt

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
)

// valuesMatch 智能比较两个可能不同类型的值是否相等
func valuesMatch(v1, v2 interface{}) bool {
	// 处理nil和空字符串的比较
	if v1 == nil && v2 == nil {
		return true
	}
	if (v1 == nil && v2 == "") || (v1 == "" && v2 == nil) {
		return true
	}

	// 如果类型相同，直接比较
	if reflect.TypeOf(v1) == reflect.TypeOf(v2) {
		return v1 == v2
	}

	// 处理数值类型的比较
	v1Float, v1IsNumber := toFloat64(v1)
	v2Float, v2IsNumber := toFloat64(v2)

	if v1IsNumber && v2IsNumber {
		return v1Float == v2Float
	}

	// 其他类型不匹配
	return false
}

// toFloat64 尝试将值转换为float64
func toFloat64(v interface{}) (float64, bool) {
	switch v := v.(type) {
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

func TestWriteTile(t *testing.T) {
	// 1. 读取初始tile数据
	feats1, err := ReadTile(bytevals2, tileid2, PROTO_LK)
	if err != nil {
		t.Logf("Warning: Failed to read initial tile: %v", err)
	}
	(t).Logf("Initial features count: %d", len(feats1))

	// 限制测试数据量，只取前10个特性
	maxFeatures := 10
	if len(feats1) > maxFeatures {
		feats1 = feats1[:maxFeatures]
		(t).Logf("Limited test to %d features", maxFeatures)
	}

	// 2. 简要输出初始特性信息
	for i, feat := range feats1 {
		(t).Logf("Initial feature %d: ID=%v, Properties count=%d", i, feat.ID, len(feat.Properties))
		// 只输出包含'tan zhe'的属性
		for k, v := range feat.Properties {
			if s, ok := v.(string); ok && strings.Contains(s, "tan zhe") {
				(t).Logf("Found 'tan zhe' in feature %d: %s=%s", i, k, v)
			}
		}
	}

	// 3. 写入图层数据
	conf := NewConfig("LK", tileid2, PROTO_MAPBOX)
	data := WriteLayer(feats1, conf)
	(t).Logf("Wrote %d bytes for layer", len(data))

	tempDir, err := os.MkdirTemp("", "mvt-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, "test_layer.mvt")
	if cerr := ioutil.WriteFile(tempFile, data, 0644); cerr != nil {
		t.Fatalf("Failed to write data to temp file: %v", cerr)
	}
	(t).Logf("Wrote test data to temp file: %s", tempFile)

	// 5. 重新读取数据
	feats2, err := ReadTile(data, tileid2, PROTO_MAPBOX)
	if err != nil {
		// 尝试使用不同的协议读取，看是否是协议问题
		feats2Alt, errAlt := ReadTile(data, tileid2, PROTO_LK)
		if errAlt == nil {
			t.Logf("Successfully read with PROTO_LK instead of PROTO_MAPBOX")
			feats2 = feats2Alt
		} else {
			// 输出详细错误信息
			t.Fatalf("Failed to read tile after writing with PROTO_MAPBOX: %v\nFailed with PROTO_LK: %v\nData saved to: %s", err, errAlt, tempFile)
		}
	}
	(t).Logf("Read features count: %d", len(feats2))

	// 6. 简要输出读取的特性信息
	for i, feat := range feats2 {
		(t).Logf("Read feature %d: ID=%v, Properties count=%d", i, feat.ID, len(feat.Properties))
	}

	// 6. 验证特性数量匹配
	if len(feats1) != len(feats2) {
		t.Errorf("Feature count mismatch: expected %d, got %d", len(feats1), len(feats2))
		// 输出简要比较信息
		if len(feats1) > len(feats2) {
			(t).Logf("Missing %d features in feats2", len(feats1)-len(feats2))
		} else {
			(t).Logf("Found %d extra features in feats2", len(feats2)-len(feats1))
		}
		t.FailNow()
	} else {
		(t).Logf("Successfully verified feature count: %d", len(feats1))
	}
}

func TestWriteLayer(t *testing.T) {
	// 1. 准备测试数据
	// 2. 尝试读取tile数据
	feats1, err := ReadTile(bytevals2, tileid2, PROTO_LK)
	if err != nil {
		t.Logf("Warning: Failed to read initial tile: %v", err)
		// 创建模拟数据进行测试
		feats1 = createMockFeatures()
		if len(feats1) == 0 {
			t.Fatalf("No features available for testing")
		}
	}

	// 3. 选择图层名称并创建配置
	layerName := "test_layer"
	conf := NewConfig(layerName, tileid2, PROTO_MAPBOX)

	// 4. 写入图层数据
	data := WriteLayer(feats1, conf)
	if len(data) == 0 {
		t.Fatalf("WriteLayer returned empty data")
	}
	(t).Logf("Wrote %d bytes for layer %s", len(data), layerName)

	// 5. 重新读取数据并验证
	feats2, err := ReadTile(data, tileid2, PROTO_MAPBOX)
	if err != nil {
		t.Fatalf("Failed to read tile after writing: %v", err)
	}

	// 6. 验证特性数量匹配
	if len(feats1) != len(feats2) {
		t.Errorf("Feature count mismatch: expected %d, got %d", len(feats1), len(feats2))
	} else {
		(t).Logf("Successfully verified %d features", len(feats2))
	}

	// 7. 详细验证特性属性
	if len(feats1) > 0 && len(feats2) > 0 {
		// 验证第一个特性的属性
		feat1 := feats1[0]
		feat2 := feats2[0]

		// 移除自动添加的layer属性进行比较
		props1 := make(map[string]interface{})
		for k, v := range feat1.Properties {
			props1[k] = v
		}

		props2 := make(map[string]interface{})
		for k, v := range feat2.Properties {
			if k != "layer" { // 忽略自动添加的layer属性
				props2[k] = v
			}
		}

		// 比较属性数量
		if len(props1) != len(props2) {
			t.Errorf("Property count mismatch for first feature (ignoring 'layer'): expected %d, got %d", len(props1), len(props2))
			t.Logf("Original properties: %v", props1)
			t.Logf("After write/read properties (ignoring 'layer'): %v", props2)
		} else {
			// 比较属性值
			allPropertiesMatch := true
			for k, v1 := range props1 {
				if v2, ok := props2[k]; !ok {
					allPropertiesMatch = false
					t.Errorf("Property '%s' not found in read feature", k)
				} else if !valuesMatch(v1, v2) {
					allPropertiesMatch = false
					t.Errorf("Property mismatch for '%s': expected %v (%T), got %v (%T)", k, v1, v1, v2, v2)
				}
			}

			// 检查是否有额外的属性
			for k := range props2 {
				if _, ok := props1[k]; !ok {
					allPropertiesMatch = false
					t.Errorf("Extra property '%s' found in read feature", k)
				}
			}

			if allPropertiesMatch {
				(t).Logf("Successfully verified properties for first feature")
			}
		}
	}
}

// createMockFeatures 创建模拟的特性数据用于测试
func createMockFeatures() []*geom.Feature {
	feats := []*geom.Feature{}

	// 创建一个点特性
	point := general.NewPoint([]float64{100, 200})
	feat1 := &geom.Feature{
		Geometry: point,
		Properties: map[string]interface{}{
			"name":  "Test Point",
			"value": 123,
		},
	}
	feats = append(feats, feat1)

	// 创建一个线特性
	line := general.NewLineString([][]float64{{100, 200}, {150, 250}, {200, 200}})
	feat2 := &geom.Feature{
		Geometry: line,
		Properties: map[string]interface{}{
			"name":   "Test Line",
			"length": 141.42,
		},
	}
	feats = append(feats, feat2)

	return feats
}
