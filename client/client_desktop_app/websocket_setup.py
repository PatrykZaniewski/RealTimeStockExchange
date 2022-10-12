import datetime
import json
import os

from PyQt5 import QtCore, QtWebSockets
from PyQt5.QtCore import QUrl

from client.client_desktop_app.gui import get_main_window, MainWindow
from client.client_desktop_app.model.order_status import OrderStatus


class QClient(QtCore.QObject):
    def __init__(self, parent):
        super().__init__(parent)

        self.client = QtWebSockets.QWebSocket()
        self.client.textMessageReceived.connect(self.on_message)
        self.client.connected.connect(self.on_connected)

        self.client.open(QUrl(os.getenv("data_streamer_websocket")))

    def on_connected(self):
        pass
        self.client.sendBinaryMessage(str.encode(os.getenv("identifier")))

    def on_message(self, message):
        message = json.loads(message)
        if message['Type'] == "PRICE":
            print(message)
            main_window: MainWindow = get_main_window()
            main_window.update_price(message)
        elif message['Type'] == 'ORDER_STATUS':
            # status = OrderStatus(**message)
            print(f"{message['Id']},RECEIVED,{datetime.datetime.timestamp(datetime.datetime.now())}")
        self.client.ping(b"ping")

    def close(self):
        self.client.close()