/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-11 10:14
* Description:
*****************************************************************/

package main

import (
	"flag"
	"fmt"
	"github.com/go-xe2/x/os/xfile"
	builder "github.com/go-xe2/xthrift/builder"
	"github.com/go-xe2/xthrift/builder/gcontext"
	"github.com/go-xe2/xthrift/pdl"
)


var (
	argDirName = flag.String("d", "", "输入yaml|json协议项目路径")
	argOutputType = flag.String("t", "", "输入文件类型,支持输出格式json|yaml|thrift" )
	argOutputDir = flag.String("o", "", "输出目录")
	argHelp = flag.Bool("h", false, "显示帮助")
	argLang = flag.String("l", "go", "生成thrift协议中命名空间的语言, 多个语言使用逗号分割")
	argGen = flag.String("gen", "", "生成目标代码,支持goLang")
	argGenArg = flag.String("arg", "", "生成目标代码参数, 目标代码为golang时，arg为模块名前掇")
	argProjName = flag.String("n", "", "设置项目名称")
	argBuild = flag.Bool("build", false, "打包协议项目文件")
)

func getDirFiles(dir string) []string {
	files1, err := xfile.ScanDir(dir, "*.json", true)
	if err != nil {
		panic(err)
	}
	files2, err := xfile.ScanDir(dir, "*.yaml", true)
	if err != nil {
		panic(err)
	}
	files1 = append(files1, files2...)
	return files1
}

func main()  {
	flag.Parse()
	if *argHelp || *argDirName == ""{
		flag.Usage()
		return
	}

	defer func() {
		if e := recover(); e != nil {
			fmt.Println("错误:", e)
		}
	}()

	project, err := pdl.NewFileProject(*argDirName)
	if err != nil {
		panic(err)
	}

	if err := project.Load(); err != nil {
		panic(err)
	}

	//if err := project.Check(); err != nil {
	//	panic(err)
	//}

	if *argProjName != "" {
		project.SetProjectName(*argProjName)
	}

	errs := project.Errors()
	for _, err := range errs {
		fmt.Println("错误:" + err.Error())
	}

	outType := ""
	if *argOutputType == "json" {
		outType = "json"
	} else if *argOutputType == "yaml" {
		outType = "yaml"
	} else if *argOutputType == "thrift" {
		outType = "thrift"
	}
	if outType == "" && *argGen == "" && ! *argBuild {
		fmt.Println("所有协议已经检查完成")
		return
	}

	outDir := "." + xfile.Separator
	if *argOutputDir != "" {
		outDir = *argOutputDir
		if !xfile.Exists(outDir) {
			err := xfile.Mkdir(outDir)
			if err != nil {
				panic(err)
			}
		}
	}


	// 生成目标代码
	if *argGen != "" {
		if err := genLangCode(outDir, *argGen, *argGenArg, project); err != nil {
			panic(err)
		}
		fmt.Printf("生成%s代码成功，存放目录:%s\n", *argGen, outDir)
		return
	}

	if *argBuild {
		if project.GetProjectName() == "" {
			fmt.Printf("未设置项目名称，请使用-n参数设置项目名称\n")
			return
		}
		fileName := xfile.Join(outDir, fmt.Sprintf("%s.pdl", project.GetProjectName()))
		if err := project.SaveToFile(fileName); err != nil {
			fmt.Printf("打包项目文件出错:%s\n", err.Error())
		}
		fmt.Printf("打包协议文件成功，存放路径:%s",fileName)
		return
	}

	if outType != "" {
		if err := export(outDir, outType, project, *argLang); err != nil {
			panic(err)
		} else {
			fmt.Printf("项目成功导出为%s文件, 存放目录%s\n", outType, outDir)
		}
		return
	}
}

func export(dir string, outType string, proj *pdl.FileProject, lang string) error  {
	var ioMgr = pdl.NewFileIOFileManager()

	var export pdl.FileExport

	if outType == "json" {
		export = pdl.NewFileExportJson(dir, ioMgr)
	} else if outType == "thrift"{
		if lang == "" {
			lang = "go,java,php"
		}
		export = pdl.NewFileExportThrift(dir, ioMgr, lang)
	} else {
		export = pdl.NewFileExportYaml(dir, ioMgr)
	}
	if err := proj.Export(export); err != nil {
		panic(err)
	}
	return nil
}


func genLangCode(dir string, lang string, arg string, proj *pdl.FileProject) error  {
	if lang == "golang" {
		return genGoLangCode(dir, arg, proj)
	}
	return nil
}

func genGoLangCode(dir string, arg string, proj *pdl.FileProject) error {
	cxt := gcontext.NewContext(dir, arg)
	bd := builder.NewProtoBuilder(cxt, proj)
	_, err := bd.Build()
	if err != nil {
		fmt.Println("生成代码出错:", err)
	}
	return nil
}