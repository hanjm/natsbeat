from packetbeat import BaseTest


class Test(BaseTest):

    def test_http_sample(self):
        self.render_config_template(http_ports=['8000'])
        self.run_packetbeat(pcap="http_10_connection_close.pcap",
                            debug_selectors=["http"])
        objs = self.read_output()

        assert len(objs) == 1
        obj = objs[0]

        assert obj["status"] == "OK"
        assert obj["http.response.headers"]["content-length"] == 11422
        assert obj["http.response.code"] == 200
        assert obj["type"] == "http"
        assert obj["client_ip"] == "127.0.0.1"
        assert obj["client_port"] == 37885
        assert obj["ip"] == "127.0.0.1"
        assert obj["port"] == 8000
