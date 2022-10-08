import sys

from PyQt5.QtWidgets import QApplication

from client.client_desktop_app.gui import MainWindow
from client.client_desktop_app.websocket_tmp import QClient

window: MainWindow

if __name__ == "__main__":
    app = QApplication([])
    app.setQuitOnLastWindowClosed(True)
    window = MainWindow()
    qclient = QClient(app)
    ret = app.exec_()
    sys.exit(ret)
