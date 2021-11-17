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
 * Variables de programa
 */
unsigned long timeSinceLastReport = 0; // Tiempo del ultimo envio al servidor
double ampSampleSum = 0;               // Acumulador de muestreos de amperaje
int samplesTaken = 0;                  // Cantidad de muestras acumuladas de amperaje

/**
 * Conexion a nuestro shield de internet
 */
EthernetClient ethernet;

/**
 * Instancia de cliente http con una conexion a nuestra API
 */
HttpClient client = HttpClient(ethernet, server, port);

/**
 *  Instancia de libreria para medir el consumo
 */
EnergyMonitor emon1;

/**
 * Codigo a ejecutarse al inicio
 */
void setup()
{
  // Check si la configuracion de DHCP fallo
  if (Ethernet.begin(mac) == 0)
  {
    Serial.println(F("Fallo al configurar Ethernet usando DHCP"));
    // Check si tenemos el hardware
    while (Ethernet.hardwareStatus() == EthernetNoHardware)
    {
      Serial.println(F("El ethernet shield no esta conectado!"));
      delay(1000);
    }
    // Check si el cable esta conectado
    while (Ethernet.linkStatus() == LinkOFF)
    {
      Serial.println(F("El cable de lan no esta conectado!"));
      delay(1000);
    }
  }
  // Seteamos el pin del sensor de amperaje
  emon1.current(1, 96.05);
  // esperamos a que el ethernet shield inicialice
  delay(1000);
}

/**
 * Loop a ejecutarse indefinidamente luego de setup
 */
void loop()
{
  /**
   * Acumulamos los muestreos de amperaje por 5 segundos y aumentamos nuestro contador de muestreo
   */
  ampSampleSum += emon1.calcIrms(1);
  samplesTaken += 1;

  /**
   * Enviamos un reporte en intervalos de 5 seg
   * 
   * Check si pasaron 5 segundos del ultimo envio.
   */
  if ((millis() - timeSinceLastReport) >= 5000)
  {
    /**
     * Dividimos la suma de muestreo por la cantidad de muestras tomadas para conseguir un 
     * valor mas consistente de medida.
     * 
     * Calculamos watts consumidos asumiendo una constante de 220v
     */
    double ampsMeasured = ampSampleSum / samplesTaken;
    double wattsMeasured = ampsMeasured * lineVoltage;

    /**
     * Preparamos el reporte de consumo
     */
    const int capacity = JSON_OBJECT_SIZE(4);
    StaticJsonDocument<capacity> doc;
    doc["meter_id"] = uniqueDeviceId;
    doc["date"] = "2021-11-16T12:59:20.0268798-03:00"; // <- esta fecha se remplaza en el lado del servidor, esta aca hasta que agreguemos algun dispositivo para el tiempo
    doc["watt_hour"] = wattsMeasured;
    doc["amps_hour"] = ampsMeasured;

    /**
     * Serializamos a un string JSON
     */
    size_t reportLength = measureJson(doc) + 1; // + 1 para el incluir el terminador nulo
    char report[reportLength];
    serializeJson(doc, report, reportLength);

    /**
     * Enviamos el reporte al servidor en formato json
     */
    client.beginRequest();
    client.post(reportPath);
    client.sendHeader("Authorization", uniqueDeviceSecret);
    client.sendHeader("Content-Type", "application/json");
    client.sendHeader("Content-Length", reportLength);
    client.beginBody();
    client.print(report);
    client.endRequest();

    /**
     * Guardamos el ultimo tiempo de envio y reseteamos nuestras variables
     * de muestreo 
     */
    timeSinceLastReport = millis();
    ampSampleSum = 0;
    samplesTaken = 0;
  }
}