import sys

from PyQt5.QtWidgets import QApplication
from dotenv import load_dotenv

from client.client_desktop_app.gui import MainWindow
from client.client_desktop_app.websocket_setup import QClient


if __name__ == "__main__":
    app = QApplication([])
    app.setQuitOnLastWindowClosed(True)
    window = MainWindow()
    load_dotenv(dotenv_path="./settings.env")
    qclient = QClient(app)
    ret = app.exec_()
    sys.exit(ret)
