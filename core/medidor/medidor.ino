#include <Ethernet.h>
#include <ArduinoHttpClient.h>
#include <ArduinoJson.h>
#include <EmonLib.h>
#include "secret.h"

/**
 * Constantes de programa
 */
const double lineVoltage = 220.0; // voltaje de lineas de py

/**
 * Conexion a nuestro shield de internet
 */
EthernetClient ethernet;

/**
 * Instancia de cliente http con una conexion a nuestra API
 */
HttpClient client = HttpClient(ethernet, server, port);

/**
 *  Instancia de libreria para medir el consumo de potencia
 */
EnergyMonitor emon1; // Create an instance

/**
 * Codigo a ejecutarse al inicio
 */
void setup()
{
  // Conectamos el puerto serial para debbugin
  Serial.begin(9600);
  // Seteamos el pin del sensor y voltaje de circuito a medir (220v paraguay)
  emon1.current(1, 80); // Current: input pin, calibration.

  if (Ethernet.begin(mac) == 0)
  {
    Serial.println(F("Fallo al configurar Ethernet usando DHCP"));
    // Check for Ethernet hardware present
    if (Ethernet.hardwareStatus() == EthernetNoHardware)
    {
      Serial.println(F("El ethernet shield no esta conectado!"));
      while (true)
      {
        delay(1);
      }
    }
    if (Ethernet.linkStatus() == LinkOFF)
      Serial.println(F("El cable de lan no esta conectado!"));
  }
  else
  {
    Serial.print(F("  DHCP asigno IP "));
    Serial.println(Ethernet.localIP());
  }
  // esperamos a que el ethernet shield inicialice
  delay(1000);

  Serial.println(F("Conectando al servidor..."));
}

/**
 * Tiempo del ultimo envio al servidor
 */
unsigned long timeSinceLastReport = 0;

double amps = 0;
long samples = 0;

/**
 * Loop a ejecutarse indefinidamente luego de setup
 */
void loop()
{
  amps = amps + emon1.calcIrms(4000);
  samples = samples + 1;
  // Check si paso 5 segundos del ultimo envio
  if ((millis() - timeSinceLastReport) >= 5000)
  {
    double mamps = (amps / samples);
    double watts = (mamps * lineVoltage);
    Serial.print("Samples: ");
    Serial.println(samples);
    Serial.print("Amps: ");
    Serial.println(mamps);
    Serial.print("Watts: ");
    Serial.println(watts);
    /**
   * Preparamos el reporte de consumo
   */
    const int capacity = JSON_OBJECT_SIZE(4);
    StaticJsonDocument<capacity> doc;
    doc["meter_id"] = uniqueDeviceId;
    doc["date"] = "2021-11-16T12:59:20.0268798-03:00"; // <- esta fecha se remplaza en el lado del servidor, esta aca hasta que agreguemos algun dispositivo para el tiempo
    doc["watt_hour"] = watts;
    doc["amps_hour"] = mamps;

    /**
   * Serializamos a un string JSON
   */
    size_t reportLength = measureJson(doc) + 1; // + 1 para el incluir el terminador
    char report[reportLength];
    serializeJson(doc, report, reportLength);

    /**
    * Enviamos el reporte al servidor 
    */
    Serial.println("Sending report...");
    client.beginRequest();
    client.post(reportPath);
    client.sendHeader("Authorization", uniqueDeviceSecret);
    client.sendHeader("Content-Type", "application/json");
    client.sendHeader("Content-Length", reportLength);
    client.beginBody();
    client.print(report);
    client.endRequest();

    // Guardamos el ultimo tiempo de envio
    timeSinceLastReport = millis();
    amps = 0;
    samples = 0;
  }
}