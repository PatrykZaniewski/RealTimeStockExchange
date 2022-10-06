import asyncio
import sys
import traceback
from asyncio import coroutine

import rel as rel
import websocket as websocket
from PyQt5.QtCore import QTimer, QObject, pyqtSignal, QRunnable, pyqtSlot, QThreadPool
from PyQt5.QtWidgets import QApplication, QVBoxLayout, QWidget, QLabel, QPushButton, QMainWindow

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
    websocket.enableTrace(True)
    list = Listener(tekst)
    ws = websocket.WebSocketApp("ws://localhost:5014/ws",
                                on_open=list.on_open,
                                on_message=list.on_message,
                                on_error=list.on_error,
                                on_close=list.on_close)

    ws.run_forever()  # Set dispatcher to automatic reconnection
    # rel.signal(2, rel.abort)  # Keyboard Interrupt
    # rel.dispatch()

class WorkerSignals(QObject):
    finished = pyqtSignal()
    error = pyqtSignal(tuple)
    result = pyqtSignal(object)
    progress = pyqtSignal(int)


class Worker(QRunnable):
    def __init__(self, fn, *args, **kwargs):
        super(Worker, self).__init__()

        # Store constructor arguments (re-used for processing)
        self.fn = fn
        self.args = args
        self.kwargs = kwargs
        self.signals = WorkerSignals()

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

        self.counter = 0

        layout = QVBoxLayout()

        self.l = QLabel("Start")
        self.l.setText("XD")
        b = QPushButton("DANGER!")
        b.pressed.connect(self.oh_no)

        layout.addWidget(self.l)
        layout.addWidget(b)

        w = QWidget()
        w.setLayout(layout)

        self.setFixedWidth(500)
        self.setFixedHeight(500)

        self.setCentralWidget(w)

        self.show()

        self.threadpool = QThreadPool()

        self.timer = QTimer()
        self.timer.setInterval(1000)
        self.timer.timeout.connect(self.recurring_timer)
        self.timer.start()

    def oh_no(self):
        # Pass the function to execute
        worker = Worker(hello, tekst=self.l)  # Any other args, kwargs are passed to the run function
        # Execute
        self.threadpool.start(worker)

    def recurring_timer(self):
        self.counter += 1
        self.l.setText("Counter: %d" % self.counter)


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
