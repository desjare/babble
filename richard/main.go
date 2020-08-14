//source: http://doc.qt.io/qt-5/qtwidgets-richtext-textedit-example.html

package main

import (
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)


func main() {
	widgets.NewQApplication(len(os.Args), os.Args)

	//create a window
	var window = widgets.NewQMainWindow(nil, 0)
	window.SetWindowTitle("Richard")
	window.SetMinimumSize2(200, 200)

	//create a layout
	var layout = widgets.NewQVBoxLayout()

	//add the layout to the centralWidget
	var centralWidget = widgets.NewQWidget(window, 0)
	centralWidget.SetLayout(layout)

	textEdit := widgets.NewQTextEdit(nil)

	//add the button to the layout
	layout.AddWidget(textEdit, 0, core.Qt__AlignCenter)

	//add the centralWidget to the window and show the window
	window.SetCentralWidget(centralWidget)
	window.Show()

	widgets.QApplication_Exec()

}
