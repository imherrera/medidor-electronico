# Diagrama de flujo general
![image info](https://firebasestorage.googleapis.com/v0/b/pandora-db134.appspot.com/o/delete_later%2Fmedidor_electronico_online.drawio%20(1).png?alt=media&token=8c88ac9a-82ce-430a-8020-12fea45ca486)

# Requisitos
###### Hardware:
* 1 x Arduino
* 1 x Shield Ethernet W5100
* 1 x Sensor CT YHDC SCT-013
* 1 x Resistencia de 18 Ohms si se utilizara 3.3V, o 33 Ohms para 5V
* 2 x 10k Ohm resistores
* 1 x 10uF capacitores

# Estructura del proyecto
###### Server:
Este directorio contiene todo el codigo relacionado a el backend

###### Client:
Este directorio contiene codigo relacionado a la/s plataformas de visualizacion

###### Core: 
Este directorio contiene el codigo relacionado al medidor

# Lenguajes, frameworks y uso
* Go -> Servidor
* JS, HTML, CSS: React, ANTD -> Pagina Web
* Kotlin: Compose, Ktor, Serialization -> App Android
* C++ -> Arduino
