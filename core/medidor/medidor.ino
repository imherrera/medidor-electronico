#include <Ethernet.h>
#include <ArduinoHttpClient.h>
#include <ArduinoJson.h>
#include <EmonLib.h>

/**
 * Constantes de programa
 */
const char uniqueDeviceId[] = "41250050123";       // id unico para identificar este medidor en el servidor
const char reportPath[] = "/report/watt-hour";     // ruta para reporte de consumo
const char server[] = "192.168.0.12";              // host de la api
const short port = 8080;                           // puerto de acceso
byte mac[] = {0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xED}; // Direccion MAC

/**
 * Conexion a nuestro shield de internet
 */
EthernetClient ethernet;

/**
 * Instancia de cliente http con una conexion a nuestra API
 */
HttpClient client = HttpClient(ethernet, server, port);

/**
 * Codigo a ejecutarse al inicio
 */
void setup()
{
  Serial.begin(9600);

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
 * Loop a ejecutarse indefinidamente luego de setup
 */
void loop()
{
  // TODO: remove
  float usage = random(20, 50);
  Serial.print("Usage: ");
  Serial.println(usage);

  /**
   * Preparamos el reporte de consumo
   */
  const int capacity = JSON_OBJECT_SIZE(3);
  StaticJsonDocument<capacity> doc;
  doc["meter_id"] = uniqueDeviceId;
  doc["date"] = "2021-11-16T12:59:20.0268798-03:00";
  doc["watt_hour"] = usage;

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
  client.sendHeader("Content-Type", "application/json");
  client.sendHeader("Content-Length", reportLength);
  client.beginBody();
  client.print(report);
  client.endRequest();

  Serial.println("Wait five seconds");
  delay(5000);
}