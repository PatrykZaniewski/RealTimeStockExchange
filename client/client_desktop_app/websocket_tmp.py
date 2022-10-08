import json

from PyQt5 import QtCore, QtWebSockets
from PyQt5.QtCore import QUrl

from client.client_desktop_app.gui import get_main_window, MainWindow


class QClient(QtCore.QObject):
    def __init__(self, parent):
        super().__init__(parent)

        self.client = QtWebSockets.QWebSocket()
        self.client.textMessageReceived.connect(self.on_message)

        self.client.open(QUrl("ws://localhost:5014/ws"))

    def on_message(self, message):
        message = json.loads(message)
        print(message)
        main_window: MainWindow = get_main_window()
        main_window.update_price(message)
        self.client.ping(b"ping")

    def close(self):
        self.client.close()