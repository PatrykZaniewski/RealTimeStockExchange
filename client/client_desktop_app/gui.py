import traceback
from typing import Optional

from PyQt5.QtCore import QRunnable, pyqtSlot
from PyQt5.QtGui import QIntValidator
from PyQt5.QtWidgets import QWidget, QHBoxLayout, QVBoxLayout, QLabel, QPushButton, QMainWindow, QLineEdit, QApplication

from client.client_desktop_app.business_logic import process_order, OrderType


def get_main_window():
    for widget in QApplication.topLevelWidgets():
        if isinstance(widget, QMainWindow):
            return widget
    return None


class QWorker(QRunnable):
    def __init__(self, fn, *args, **kwargs):
        super(QWorker, self).__init__()

        self.fn = fn
        self.args = args
        self.kwargs = kwargs

    @pyqtSlot()
    def run(self):
        try:
            self.fn(*self.args, **self.kwargs)
        except Exception as _:
            traceback.print_exc()


class MainWindow(QMainWindow):
    def __init__(self, *args, **kwargs):
        super(MainWindow, self).__init__(*args, **kwargs)
        self.setFixedWidth(400)
        self.setFixedHeight(400)
        self.setWindowTitle("Broker Client")

        v_layout = QVBoxLayout()

        for asset in ["ASSECO", "COMARCH", "CDPROJECT"]:
            v_layout.addLayout(self.init_assets(asset))

        w = QWidget()
        w.setLayout(v_layout)
        self.v_layout = v_layout
        self.setCentralWidget(w)

        self.show()

        # self.threadpool = QThreadPool()
        # worker = QWorker(init_websocket_connection)
        # self.threadpool.start(worker)

    def init_assets(self, asset_name: str) -> QHBoxLayout:
        h_layout = QHBoxLayout()
        h_layout.setObjectName(asset_name)

        label = QLabel(asset_name)
        buy_price = QLabel()
        buy_price.setText("100.00")
        buy_price.setObjectName("buy_price")
        sell_price = QLabel()
        sell_price.setText("50.00")
        sell_price.setObjectName("sell_price")
        validator = QIntValidator()
        validator.setRange(0, 100)
        sell_amount = QLineEdit()
        sell_amount.setValidator(validator)
        buy_amount = QLineEdit()
        buy_amount.setValidator(validator)
        buy = QPushButton("BUY!")
        buy.pressed.connect(lambda: self._process_order(asset_name, buy_amount, buy_price, OrderType.BUY))
        sell = QPushButton("SELL!")
        sell.pressed.connect(lambda: self._process_order(asset_name, sell_amount, sell_price, OrderType.SELL))
        h_layout.addWidget(label)
        h_layout.addWidget(buy_price)
        h_layout.addWidget(sell_price)
        h_layout.addWidget(buy_amount)
        h_layout.addWidget(sell_amount)
        h_layout.addWidget(buy)
        h_layout.addWidget(sell)

        return h_layout

    def update_price(self, asset_data):
        layout: Optional[QHBoxLayout] = None
        for i in range(self.v_layout.count()):
            if self.v_layout.itemAt(i).layout().objectName() == asset_data.get("AssetName"):
                layout = self.v_layout.itemAt(i).layout()
                break

        for i in range(layout.count()):
            if layout.itemAt(i).widget().objectName() == "buy_price":
                layout.itemAt(i).widget().setText(str(asset_data.get("BuyPrice")))
            if layout.itemAt(i).widget().objectName() == "sell_price":
                layout.itemAt(i).widget().setText(str(asset_data.get("SellPrice")))

    def _process_order(self, asset_name: str, amount: QLineEdit, price: QLabel, order_type: OrderType):
        process_order(asset_name, int(amount.text()), float(price.text()), order_type)
