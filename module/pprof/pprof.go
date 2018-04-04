package pprof

import (
	"fmt"
	"log"
	//"net/http"
	netpprof "net/http/pprof"
	"os"
	"path"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"
	"vectors/logger"
	"vectors/web"
)

var (
	PprofModule *web.TModule
	pid         int
)

func init() {
	pid = os.Getpid()

	PprofModule = web.NewModule(nil, "")
	PprofModule.Get("/debug/pprof/", netpprof.Index)
	PprofModule.Get("/debug/pprof/cmdline", netpprof.Cmdline)
	PprofModule.Get("/debug/pprof/profile", netpprof.Profile)
	PprofModule.Get("/debug/pprof/symbol", netpprof.Symbol)

	PprofModule.Get("/debug/pprof/block", block)
	PprofModule.Get("/debug/pprof/goroutine", goroutine)
	PprofModule.Get("/debug/pprof/heap", heap)
	PprofModule.Get("/debug/pprof/threadcreate", threadcreate)

	PprofModule.Get("/debug/pprof/cpu-profile", CPUProfile)
	PprofModule.Get("/debug/pprof/mem-profile", memProf)

}

func block(hd *web.THandler) {
	p := pprof.Lookup("block")
	p.WriteTo(hd, 2)
}

func goroutine(hd *web.THandler) {
	p := pprof.Lookup("goroutine")
	p.WriteTo(hd, 2)
}

func heap(hd *web.THandler) {
	p := pprof.Lookup("heap")
	p.WriteTo(hd, 2)
}

func threadcreate(hd *web.THandler) {
	p := pprof.Lookup("threadcreate")
	p.WriteTo(hd, 2)
}

// record memory profile in pprof
func memProf(hd *web.THandler) {
	filename := "mem-" + strconv.Itoa(pid) + ".mprof"
	if f, err := os.Create(filename); err != nil {
		fmt.Fprintf(hd, "create file %s error %s\n", filename, err.Error())
		log.Fatal("record heap profile failed: ", err)
	} else {
		runtime.GC()
		pprof.WriteHeapProfile(f)
		f.Close()
		fmt.Fprintf(hd, "create heap profile %s \n", filename)
		_, fl := path.Split(os.Args[0])
		fmt.Fprintf(hd, "Now you can use this to check it: go tool pprof %s %s\n", fl, filename)
	}
}

// start cpu profile monitor
func CPUProfile(hd *web.THandler) {
	// 创建pprof文件
	filename := "cpu-" + strconv.Itoa(pid) + ".pprof"
	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(hd, "Could not Creat file %s: %s\n", filename, err)
		log.Fatal("record cpu profile failed: ", err)
	}
	defer f.Close()

	//pprof.StartCPUProfile(f)
	// Set Content Type assuming StartCPUProfile will work,
	// because if it does it starts writing.
	//hd.Header().Set("Content-Type", "application/octet-stream")

	// 开始记录系统运行过程保存到<f>文件
	if err := pprof.StartCPUProfile(f); err != nil {
		// StartCPUProfile failed, so no writes yet.
		// Can change header back to text content
		// and send error code.
		//hd.Header().Set("Content-Type", "text/plain; charset=utf-8")
		//hd.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(hd, "Could not enable CPU profiling: %s\n", err)
		return
	}
	lParams := hd.MethodParams()
	// 延迟实际记录时间必须达到30秒 否则获取不到具体信息
	lSec := lParams.AsInteger("seconds")
	if lSec == 0 {
		lSec = 120
	}
	logger.Info(lSec, lParams.AsInteger("seconds"))
	time.Sleep(time.Duration(lSec) * time.Second)

	// 结束
	pprof.StopCPUProfile()

	fmt.Fprintf(hd, "create cpu profile %s \n", filename)
	root, _ := os.Getwd() //path.Split(os.Args[0])
	fmt.Fprintf(hd, "Now you can use this to check it: go tool pprof %s\n", path.Join("file:///", root, filename))
}
