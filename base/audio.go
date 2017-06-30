package base

/*
#cgo pkg-config: gstreamer-1.0
#include <gst/gst.h>


// ******************** 定义消息处理函数 ********************
gboolean bus_call(GstBus *bus, GstMessage *msg, gpointer data)
{
	GMainLoop *loop = (GMainLoop *)data;//这个是主循环的指针，在接受EOS消息时退出循环
	gchar *debug;
	GError *error;

	switch (GST_MESSAGE_TYPE(msg)) {
	case GST_MESSAGE_EOS:
		g_main_loop_quit(loop);
		//g_print("EOF\n");
		break;
	case GST_MESSAGE_ERROR:
		gst_message_parse_error(msg,&error,&debug);
		g_free(debug);
		g_printerr("ERROR:%s\n",error->message);
		g_error_free(error);
		g_main_loop_quit(loop);
		break;
	default:
		break;
	}

	return TRUE;
}

static GstBus *pipeline_get_bus(void *pipeline)
{
	return gst_pipeline_get_bus(GST_PIPELINE(pipeline));
}

static void bus_add_watch(void *bus, void *loop)
{
	gst_bus_add_watch(bus, bus_call, loop);
	gst_object_unref(bus);
}

static void set_path(void *play, gchar *path)
{
	g_object_set(G_OBJECT(play), "uri", path, NULL);
}

static void object_unref(void *pipeline)
{
	gst_object_unref(GST_OBJECT(pipeline));
}

static void media_ready(void *pipeline)
{
	gst_element_set_state(pipeline, GST_STATE_READY);
}

static void media_pause(void *pipeline)
{
	gst_element_set_state(pipeline, GST_STATE_PAUSED);
}

static void media_play(void *pipeline)
{
	gst_element_set_state(pipeline, GST_STATE_PLAYING);
}

static void media_stop(void *pipeline)
{
	gst_element_set_state(pipeline, GST_STATE_NULL);
}

static void set_mute(void *play)
{
	g_object_set(G_OBJECT(play), "mute", FALSE, NULL);
}

static void set_volume(void *play, int vol)
{
	int ret = vol % 101;

	g_object_set(G_OBJECT(play), "volume", ret/10.0, NULL);
}
static void media_seek(void *pipeline, gint64 pos)
{
	gint64 cpos;

	gst_element_query_position (pipeline, GST_FORMAT_TIME, &cpos);
	cpos += pos*1000*1000*1000;
	if (!gst_element_seek (pipeline, 1.0, GST_FORMAT_TIME, GST_SEEK_FLAG_FLUSH,
                         GST_SEEK_TYPE_SET, cpos,
                         GST_SEEK_TYPE_NONE, GST_CLOCK_TIME_NONE)) {
    		g_print ("Seek failed!\n");
    	}
}

*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"
	"strings"
	"path/filepath"
	"os"
)

var path string

func GString(s string) *C.gchar {
	return (*C.gchar)(C.CString(s))
}

func GFree(s unsafe.Pointer) {
	C.g_free(C.gpointer(s))
}

func PlayProcess(loop *C.GMainLoop) {
	var pipeline *C.GstElement // 定义组件
	var bus *C.GstBus

	v0 := GString("playbin")
	v1 := GString("play")
	pipeline = C.gst_element_factory_make(v0, v1)
	GFree(unsafe.Pointer(v0))
	GFree(unsafe.Pointer(v1))
	// 得到 管道的消息总线
	bus = C.pipeline_get_bus(unsafe.Pointer(pipeline))
	if bus == (*C.GstBus)(nil) {
		fmt.Println("GstBus element could not be created.Exiting.")
		return
	}
	C.bus_add_watch(unsafe.Pointer(bus), unsafe.Pointer(loop))

	v2 := GString(path)
	C.set_path(unsafe.Pointer(pipeline), v2)
	GFree(unsafe.Pointer(v2))

	C.media_ready(unsafe.Pointer(pipeline))
	C.media_play(unsafe.Pointer(pipeline))

	C.g_main_loop_quit(loop)
}

func main() {
	Play("b.mp3")
}
func Play(text string) {
	var loop *C.GMainLoop
	if text==""{
		fmt.Println("参数不能为空")
		return
	}
	if strings.Contains(text,"http") {
		path=text
		C.gst_init((*C.int)(unsafe.Pointer(nil)),
			(***C.char)(unsafe.Pointer(nil)))
		loop = C.g_main_loop_new((*C.GMainContext)(unsafe.Pointer(nil)),
			C.gboolean(0)) // 创建主循环，在执行 g_main_loop_run后正式开始循环
	} else  {
		_, err := os.Stat(text)
		if err != nil && os.IsNotExist(err) {
			fmt.Printf("只能是url音频或文件，请检查参数.Error: %v\n", err)
			return
		}
		p, err := filepath.Abs(text)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		path = fmt.Sprintf("file://%s", p)
		fmt.Println(path)
	}
	C.gst_init((*C.int)(unsafe.Pointer(nil)),
		(***C.char)(unsafe.Pointer(nil)))
	loop = C.g_main_loop_new((*C.GMainContext)(unsafe.Pointer(nil)),
		C.gboolean(0)) // 创建主循环，在执行 g_main_loop_run后正式开始循环

	PlayProcess(loop)
	//你猜为什么?
	time.Sleep(500 * time.Millisecond)
}
