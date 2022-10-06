import asyncio
import sys
import traceback
from asyncio import coroutine

import rel as rel
import websocket as websocket
from PyQt5.QtCore import QTimer, QObject, pyqtSignal, QRunnable, pyqtSlot, QThreadPool
from PyQt5.QtWidgets import QApplication, QVBoxLayout, QWidget, QLabel, QPushButton, QMainWindow, QHBoxLayout


class Listener:
    def __init__(self, tekst):
        self.tekst = tekst

    def on_message(self, ws, message):
        self.tekst.setText("!!!!!!!!")
        print(message)

    def on_error(self, ws, error):
        print(error)

    def on_close(self, ws, close_status_code, close_msg):
        print("### closed ###")

    def on_open(self, ws):
        print("Opened connection")


def hello(tekst):
    list = Listener(tekst)
    ws = websocket.WebSocketApp("ws://localhost:5014/ws",
                                on_open=list.on_open,
                                on_message=list.on_message,
                                on_error=list.on_error,
                                on_close=list.on_close)

    ws.run_forever()  # Set dispatcher to automatic reconnection




class Worker(QRunnable):
    def __init__(self, fn, *args, **kwargs):
        super(Worker, self).__init__()

        # Store constructor arguments (re-used for processing)
        self.fn = fn
        self.args = args
        self.kwargs = kwargs

    @pyqtSlot()
    def run(self):
        # Retrieve args/kwargs here; and fire processing using them
        try:
            result = self.fn(*self.args, **self.kwargs)
        except:
            traceback.print_exc()
            exctype, value = sys.exc_info()[:2]


class MainWindow(QMainWindow):
    def __init__(self, *args, **kwargs):
        super(MainWindow, self).__init__(*args, **kwargs)
        self.setFixedWidth(400)
        self.setFixedHeight(400)

        v_layout = QVBoxLayout()

        h_layout_asseco = QHBoxLayout()
        h_layout_asseco.setObjectName("asseco")
        h_layout_comarch = QHBoxLayout()
        h_layout_comarch.setObjectName("comarch")
        h_layout_cdproject = QHBoxLayout()
        h_layout_cdproject.setObjectName("cdproject")

        asseco_label = QLabel("ASSECO")
        asseco_buy_price = QLabel("100.00")
        asseco_sell_price = QLabel("90.00")
        asseco_buy = QPushButton("BUY!")
        asseco_sell = QPushButton("SELL!")
        h_layout_asseco.addWidget(asseco_label)
        h_layout_asseco.addWidget(asseco_buy_price)
        h_layout_asseco.addWidget(asseco_sell_price)
        h_layout_asseco.addWidget(asseco_buy)
        h_layout_asseco.addWidget(asseco_sell)

        comarch_label = QLabel("COMARCH")
        comarch_buy_price = QLabel("200.00")
        comarch_sell_price = QLabel("190.00")
        comarch_buy = QPushButton("BUY!")
        comarch_sell = QPushButton("SELL!")
        h_layout_comarch.addWidget(comarch_label)
        h_layout_comarch.addWidget(comarch_buy_price)
        h_layout_comarch.addWidget(comarch_sell_price)
        h_layout_comarch.addWidget(comarch_buy)
        h_layout_comarch.addWidget(comarch_sell)

        cdproject_label = QLabel("CD PROJECT")
        cdproject_buy_price = QLabel("300.00")
        cdproject_sell_price = QLabel("290.00")
        cdproject_buy = QPushButton("BUY!")
        cdproject_sell = QPushButton("SELL!")
        h_layout_cdproject.addWidget(cdproject_label)
        h_layout_cdproject.addWidget(cdproject_buy_price)
        h_layout_cdproject.addWidget(cdproject_sell_price)
        h_layout_cdproject.addWidget(cdproject_buy)
        h_layout_cdproject.addWidget(cdproject_sell)

        # self.l = QLabel("Start")
        # self.l.setText("XD")

        b = QPushButton("DANGER!")
        b.pressed.connect(self.oh_no)

        # v_layout.addWidget(self.l)
        # v_layout.addWidget(b)

        v_layout.addLayout(h_layout_asseco)
        v_layout.addLayout(h_layout_comarch)
        v_layout.addLayout(h_layout_cdproject)

        w = QWidget()
        w.setLayout(v_layout)

        self.setCentralWidget(w)

        self.show()

        self.threadpool = QThreadPool()

    def oh_no(self):
        worker = Worker(hello, tekst=self.l)
        self.threadpool.start(worker)


if __name__ == "__main__":
    # app = QApplication([])
    # text_area = QTextEdit()
    # text_area.setFocusPolicy(Qt.NoFocus)
    # message = QLineEdit()
    # layout = QVBoxLayout()
    # layout.addWidget(text_area)
    # layout.addWidget(message)
    # window = QWidget()
    # window.setLayout(layout)
    # window.show()
    # window.thread()
    # asyncio.run(hello())
    # app.exec()
    app = QApplication([])
    window = MainWindow()
    # task = asyncio.run(hello())
    app.exec_()

# class Window(Qt.Dialog):
#     def __init__(self, parent=None):
#         super(Window, self).__init__()
#
#         self.thread = ListenWebsocket()
#         self.thread.start()
#
#
# class ListenWebsocket(QtCore.QThread):
#     def __init__(self, parent=None):
#         super(ListenWebsocket, self).__init__(parent)
#
#         websocket.enableTrace(True)
#
#         self.WS = websocket.WebSocketApp("ws://localhost:8080/chatsocket",
#                                          on_message = self.on_message,
#                                          on_error = self.on_error,
#                                          on_close = self.on_close)
#
#     def run(self):
#         #ws.on_open = on_open
#
#         self.WS.run_forever()
#
#
#     def on_message(self, ws, message):
#         print (message)
#
#     def on_error(self, ws, error):
#         print (error)
#
#     def on_close(self, ws):
#         print ("### closed ###")
#
# if __name__ == '__main__':
#     app = QtGui.QApplication(sys.argv)
#
#     QtGui.QApplication.setQuitOnLastWindowClosed(False)
#
#     window = Window()
#     window.show()
#
#     sys.exit(app.exec_())
