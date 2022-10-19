import json
import logging
import os
import time

from PyQt5 import QtCore, QtWebSockets
from PyQt5.QtCore import QUrl

from client.client_desktop_app.gui import get_main_window, MainWindow
from client.client_desktop_app.model.order_status import OrderStatus
from client.client_desktop_app.model.price import Price


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
        if message['type'] == "PRICE":
            print(message)
            price = Price.from_dict(message)
            main_window: MainWindow = get_main_window()
            main_window.update_price(price)
        elif message['type'] == 'ORDER_STATUS':
            order_status = OrderStatus.from_dict(message)
            # status = OrderStatus(**message)
            logging.info(f'{order_status.id},CLIENT,STATUS_RECEIVED,{int(time.time() * 1000000)}')
            # print(f"{order_status.id},CLIENT,STATUS_RECEIVED,{int(time.time() * 1000000)}")
        self.client.ping(b"ping")

    def close(self):
        self.client.close()
