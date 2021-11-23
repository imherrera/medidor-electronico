import './Dashboard.css';
import { replaceLast } from '../utils';

import { Navigate } from 'react-router-dom'
import { useState, useEffect } from "react";

import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend, } from 'recharts';


const kwhHomePrice = 365.45;
//const kwhBussinessPrice = 334.798

const loadDashboard = async (jwt, uci) => {
    //await sleep(1000)
    if (!jwt || !uci) return { error: 401 }

    /**
    * Parametros del request
    * **/
    const params = {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${jwt}`,
        },
    };

    /**
     * Llamada a nuestro servidor para que nos devuelve el resumen de uso
     * correspondiente a el usuario con C.I
     * **/
    const response = await fetch(`https://medidor-electronico-server.herokuapp.com/usage/resume/${uci}`, params);

    if (!response.ok) return { error: response.status }

    return await response.json();
}
var guaranies = new Intl.NumberFormat('py-PY', {
    style: 'currency',
    currency: 'PYG',
    // These options are needed to round to whole numbers if that's what you want.
    //minimumFractionDigits: 0, // (this suffices for whole numbers, but will print 2500.10 as $2,500.1)
    //maximumFractionDigits: 0, // (causes 2500.99 to be printed as $2,501)
});


function Dashboard() {
    const [refresh, setRefresh] = useState(true);
    const [error, setError] = useState(null);
    const [data, setData] = useState(null);

    useEffect(() => {
        /**
         * Leemos el token guardado, sera nulo si el usuario no tiene una sesion abierta
         * **/
        const token = localStorage.getItem('power-meter-jwt');
        const uci = localStorage.getItem('power-meter-uci');
        /**
         * Llamamos al servidor para conseguir los datos de consumo
         * **/
        loadDashboard(token, uci).then(res => {
            if (res.error) {
                setError(res.error)
                return
            }
            console.log("data: ", res)
            /**
             * Calculo de costo
             * **/
            const getCost = (wh) => {
                // 1000 -> 1 kwatt
                const kwh = ((wh / 60) / 1000);
                return kwh * kwhHomePrice
            };

            /**
             * Mapeador de respuesta de api a objeto para representar en graficos
             * **/
            const mapToUi = (e) => ({
                yname: 'W',
                xname: '5s',
                watts: e.watt_hour,
                amps: e.amps_hour,
                pv: new Date(e.date).toTimeString(),
            });


            const graph = [];
            let consumption = 0;
            let cost = 0;
            for (let i in res.reverse()) {
                const e = res[i];
                cost += getCost(e.watt_hour);
                // suma del consumo de las ultimas 24 horas
                if (i <= 17280) {
                    consumption += e.watt_hour;
                }

                if (i <= 12) graph.push(mapToUi(e));
            }
            setData({
                graph: graph.reverse(),
                cost: cost,
                consumption: (consumption / 17280)// consumo de watts por hora
            });

            // Volvemos a hacer esta llamada luego de 5seg
            setTimeout(() => setRefresh(!refresh), 5000);
        }).catch(err => {
            console.error("Failure on getting resume:", err);
        })

    }, [refresh]);


    /**
     * Redirigimos al usuario a la pantalla de inicio de sesion
     * **/
    if (error === 401) {
        return <Navigate to='/login' />
    }

    return (
        <div className="dashboard-layout">
            <div className="month-resume" title={<h1>Resumen del mes</h1>}>
                <div><h3>Resumen de tarifas y consumo</h3></div>
                <div className="flex-container">
                    <div className="col">
                        <h5>Ahora</h5>
                        <h1>{(data && data.graph.length > 0) ? data.graph[data.graph.length - 1].watts.toFixed(2) : 'calculando...'} W</h1>
                        <h1>{(data && data.graph.length > 0) ? data.graph[data.graph.length - 1].amps.toFixed(2) : 'calculando...'} A</h1>
                    </div>
                    <div className="col">
                        <h5>Hoy</h5>
                        <h1>{(data) ? guaranies.format(data.cost) : "calculando..."}</h1>
                        <h1>{(data) ? data.consumption.toFixed(0) : 'calculando...'} Wh</h1>
                    </div>
                    <div className="col">
                        <h5>Ayer</h5>
                        <h1>0 PYG</h1>
                        <h1>0 Wh</h1>
                    </div>
                    <div className="col">
                        <h5>Este mes</h5>
                        <h1>{(data) ? guaranies.format(data.cost) : "calculando..."}</h1>
                        <h1>{(data) ? (data.consumption / 1000).toFixed(2) : 'calculando...'} kWh</h1>
                    </div>
                    <div className="col">
                        <h5>El mes pasado</h5>
                        <h1>0 PYG</h1>
                        <h1>0 kWh</h1>
                    </div>
                    <div className="col">
                        <h5>Tarifa domicilio particular por kWh</h5>
                        <h1>365.45 PYG</h1>
                    </div>
                    <div className="col">
                        <h5>Tarifa industrial por kWh</h5>
                        <h1>296.56 PYG</h1>
                    </div>
                </div>
            </div>

            <h1 className="chart-title">Consumo del ultimo minuto</h1>
            <ResponsiveContainer height={300}>
                <LineChart
                    data={(data) ? data.graph : []}
                    margin={{
                        top: 10,
                        right: 30,
                        left: 0,
                        bottom: 0,
                    }}
                >
                    <CartesianGrid strokeDasharray={10} />
                    <XAxis dataKey="xname" />
                    <YAxis />
                    <Tooltip />
                    <Line type="monotone" dataKey="watts" stroke="#8884d8" />
                    <Line type="monotone" dataKey="amps" stroke="#82ca9d" />
                    <Legend />
                </LineChart>
            </ResponsiveContainer>
        </div>
    )
}

export default Dashboard;
