import sys

from PyQt5.QtWidgets import QApplication

from client.client_desktop_app.gui import MainWindow

window: MainWindow

if __name__ == "__main__":
    app = QApplication([])
    app.setQuitOnLastWindowClosed(True)
    window = MainWindow()
    ret = app.exec_()
    sys.exit(ret)
