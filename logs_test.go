package logs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetLalamoveLoggerPassDebug(t *testing.T) {
	temp, err := ioutil.TempFile(".", "zap-prod-config-test")
	require.Nil(t, err, "Failed to create temp file.")
	defer os.Remove(temp.Name())

	Config = NewConfig()
	Config.OutputPaths = []string{temp.Name()}
	Config.EncoderConfig = *NewLalamoveEncoderConfig()
	Config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	const msg = "I am a message"

	Log, err = Config.Build()
	require.Nil(t, err, "Failed to create logger")
	Logger().Debug(msg)
	byteContents, err := ioutil.ReadAll(temp)
	defer Logger().Sync()

	var tmpJSON map[string]interface{}
	err = json.Unmarshal(byteContents, &tmpJSON)
	require.Nil(t, err, "Failed to unmarshal json")
	require.Equal(t, "debug", tmpJSON["level"])
	require.Equal(t, msg, tmpJSON["message"])
	require.Contains(t, tmpJSON["src_file"], "logs/logs_test.go")
	_, err = strconv.Atoi(tmpJSON["src_line"].(string))
	require.Nil(t, err)
	_, err = time.Parse(ISO8601, tmpJSON["time"].(string))
	require.Nil(t, err, "Invalid time format")
}

func TestGetLalamoveLoggerPassDebugWithExtraFields(t *testing.T) {
	temp, err := ioutil.TempFile(".", "zap-prod-config-test")
	require.Nil(t, err, "Failed to create temp file.")
	defer os.Remove(temp.Name())

	Config = NewConfig()
	Config.OutputPaths = []string{temp.Name()}
	Config.EncoderConfig = *NewLalamoveEncoderConfig()
	Config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	const msg = "I am a message"

	Log, err = Config.Build()
	require.Nil(t, err, "Failed to create logger")
	Logger().Debug(msg, zap.String("key1", "extra field"))
	byteContents, err := ioutil.ReadAll(temp)
	defer Logger().Sync()

	var tmpJSON map[string]interface{}
	err = json.Unmarshal(byteContents, &tmpJSON)
	require.Nil(t, err, "Failed to unmarshal json")
	require.Equal(t, "debug", tmpJSON["level"])
	require.Equal(t, "extra field", tmpJSON["context"].(map[string]interface{})["key1"])
	require.Equal(t, msg, tmpJSON["message"])
	require.Contains(t, tmpJSON["src_file"], "logs/logs_test.go")
	_, err = strconv.Atoi(tmpJSON["src_line"].(string))
	require.Nil(t, err)
	_, err = time.Parse(ISO8601, tmpJSON["time"].(string))
	require.Nil(t, err, "Invalid time format")
}

func TestGetLalamoveLoggerPassInfo(t *testing.T) {
	temp, err := ioutil.TempFile(".", "zap-prod-config-test")
	require.Nil(t, err, "Failed to create temp file.")
	defer os.Remove(temp.Name())

	Config = NewConfig()
	Config.OutputPaths = []string{temp.Name()}
	Config.EncoderConfig = *NewLalamoveEncoderConfig()
	Config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	const msg = "I am a message"

	Log, err = Config.Build()
	require.Nil(t, err, "Failed to create logger")
	Logger().Info(msg)
	byteContents, err := ioutil.ReadAll(temp)
	defer Logger().Sync()

	var tmpJSON map[string]interface{}
	err = json.Unmarshal(byteContents, &tmpJSON)
	require.Nil(t, err, "Failed to unmarshal json")
	require.Equal(t, "info", tmpJSON["level"])
	require.Equal(t, msg, tmpJSON["message"])
	require.Contains(t, tmpJSON["src_file"], "logs/logs_test.go")
	_, err = strconv.Atoi(tmpJSON["src_line"].(string))
	require.Nil(t, err)
	_, err = time.Parse(ISO8601, tmpJSON["time"].(string))
	require.Nil(t, err, "Invalid time format")
}

func TestGetLalamoveLoggerPassWarn(t *testing.T) {
	temp, err := ioutil.TempFile(".", "zap-prod-config-test")
	require.Nil(t, err, "Failed to create temp file.")
	defer os.Remove(temp.Name())

	Config = NewConfig()
	Config.OutputPaths = []string{temp.Name()}
	Config.EncoderConfig = *NewLalamoveEncoderConfig()
	Config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	const msg = "I am a message"

	Log, err = Config.Build()
	require.Nil(t, err, "Failed to create logger")
	Logger().Warn(msg)
	byteContents, err := ioutil.ReadAll(temp)
	defer Logger().Sync()

	var tmpJSON map[string]interface{}
	err = json.Unmarshal(byteContents, &tmpJSON)
	require.Nil(t, err, "Failed to unmarshal json")
	require.Equal(t, Warning, tmpJSON["level"])
	require.Equal(t, msg, tmpJSON["message"])
	require.Contains(t, tmpJSON["src_file"], "logs/logs_test.go")
	_, err = strconv.Atoi(tmpJSON["src_line"].(string))
	require.Nil(t, err)
	_, err = time.Parse(ISO8601, tmpJSON["time"].(string))
	require.Nil(t, err, "Invalid time format")
}

func TestGetLalamoveLoggerPassError(t *testing.T) {
	temp, err := ioutil.TempFile(".", "zap-prod-config-test")
	require.Nil(t, err, "Failed to create temp file.")
	defer os.Remove(temp.Name())

	Config = NewConfig()
	Config.OutputPaths = []string{temp.Name()}
	Config.EncoderConfig = *NewLalamoveEncoderConfig()
	Config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	const msg = "I am a message"

	Log, err = Config.Build()
	require.Nil(t, err, "Failed to create logger")
	Logger().Error(msg)
	byteContents, err := ioutil.ReadAll(temp)
	defer Logger().Sync()

	var tmpJSON map[string]interface{}
	err = json.Unmarshal(byteContents, &tmpJSON)
	require.Nil(t, err, "Failed to unmarshal json")
	require.Equal(t, "error", tmpJSON["level"])
	require.Equal(t, msg, tmpJSON["message"])
	require.Contains(t, tmpJSON["src_file"], "logs/logs_test.go")
	require.Contains(t, tmpJSON["backtrace"], "logs/logs_test.go")
	_, err = strconv.Atoi(tmpJSON["src_line"].(string))
	require.Nil(t, err)
	_, err = time.Parse(ISO8601, tmpJSON["time"].(string))
	require.Nil(t, err, "Invalid time format")
}

func TestGetLalamoveLoggerPassLogging(t *testing.T) {
	Logger().Debug("I am a Debug")
	Logger().Info("I am an Info")
	Logger().Warn("I am a Warn")
	Logger().Error("I am an Error")
	// It should not be called as the it will return exit code 3
	// Logger().Fatal("I am not a Fatal")
	// By default, loggers are unbuffered. However, since zap's low-level APIs allow buffering,
	// calling Sync before letting your process exit is a good habit.
	defer Logger().Sync()
	require.True(t, true)
}

// TestGetLalamoveLoggerPassErrorWithRootLevelNamespace will test the extra fields.
// The extra fields should always instead the context namespace.
// expected result
//{
//    "level": "debug",
//    "time": "2018-01-03T03:40:02.087012761Z",
//    "src_file": "logs/logs_test.go",
//    "message": "I am a Debug",
//    "src_line": "27",
//    "context": {
//        "f0": "I go to school by bus",
//        "f1": "Goodest english"
//    }
//}
func TestGetLalamoveLoggerPassErrorWithRootLevelNamespace(t *testing.T) {
	Logger().Error("I am a Debug", zap.String("f0", "I go to school by bus"), zap.String("f1", "Goodest english"))
	Logger().Error("I am a Debug", zap.String("f2", "I go to school by MTR"), zap.String("f3", "Goodest cantonese"))
	defer Logger().Sync()
	require.True(t, true)
}
