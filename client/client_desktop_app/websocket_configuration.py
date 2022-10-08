import json

import rel as rel
import websocket as websocket


class Listener:
    def __init__(self, tekst):
        self.tekst = tekst

    def on_message(self, ws, message):
        from client.client_desktop_app.gui import get_main_window, MainWindow
        message = json.loads(message)
        main_window: MainWindow = get_main_window()
        main_window.update_price(message)
        print(message)
        ws.send(data={"abc": "def"})

    def on_error(self, ws, error):
        print(error)

    def on_close(self, ws, close_status_code, close_msg):
        print("### closed ###")

    def on_open(self, ws):
        print("Opened connection")


def init_websocket_connection():
    listener = Listener("tekst")
    ws = websocket.WebSocketApp("ws://localhost:5014/ws",
                                on_open=listener.on_open,
                                on_message=listener.on_message,
                                on_error=listener.on_error,
                                on_close=listener.on_close)

    ws.run_forever()
    rel.signal(2, rel.abort)
    rel.dispatch()
