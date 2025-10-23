#include <ESP8266WiFi.h>
#include "DHT.h"

#define DHTPIN 14          // ขา D5 บน NodeMCU
#define DHTTYPE DHT22
DHT dht(DHTPIN, DHTTYPE);

const char* ssid = "SSIDNAME";
const char* password = "SSIDPASSWORD";

// เปลี่ยน IP ให้เป็นของเครื่อง Termux ที่รัน logger อยู่
const char* host = "192.168.0.xx";  
const int port = 5000;

void setup() {
  Serial.begin(115200);
  dht.begin();
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("\nWiFi connected");
}

void loop() {
  float hum = dht.readHumidity();
  float temp = dht.readTemperature();
  if (isnan(hum) || isnan(temp)) {
    Serial.println("DHT read failed");
    delay(5000);
    return;
  }

  WiFiClient client;
  if (client.connect(host, port)) {
    String postData = "temp=" + String(temp, 1) + "&hum=" + String(hum, 1);
    client.println("POST /log HTTP/1.1");
    client.println("Host: " + String(host));
    client.println("Content-Type: application/x-www-form-urlencoded");
    client.print("Content-Length: ");
    client.println(postData.length());
    client.println();
    client.print(postData);
    Serial.println("Sent -> " + postData);
  } else {
    Serial.println("Connection failed");
  }

  client.stop();
  delay(600000); // 10 นาที ยิงข้อมูลที (600000 ms)
}
