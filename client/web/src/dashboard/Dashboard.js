import './Dashboard.css';
import { replaceLast } from '../utils';

import { Navigate } from 'react-router-dom'
import { useState, useEffect } from "react";

import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, } from 'recharts';


const kwhHomePrice = 402.000;
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
    const response = await fetch(`http://192.168.0.12:8080/usage/resume/${uci}`, params);

    if (!response.ok) return { error: response.status }

    return await response.json();
}

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

            /**
             * Calculo de costo
             * **/
            const getCost = (wh) => {
                // 1000 -> 1 kwatt
                const kwh = (wh / 1000);
                return kwh * kwhHomePrice
            };

            /**
             * Mapeador de respuesta de api a objeto para representar en graficos
             * **/
            const mapToUi = (e) => ({
                name: 'W/H',
                uv: e.watt_hour,
                pv: new Date(e.date).toTimeString(),
            });


            const graph = [];
            let consumption = 0;
            let cost = 0;
            for (let i in res.reverse()) {
                const e = res[i];
                cost += getCost(e.watt_hour);
                consumption += e.watt_hour;
                if (i < 24) graph.push(mapToUi(e));
            }
            setData({
                graph: graph.reverse(),
                cost: cost,
                consumption: consumption
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
                        <h5>Hoy</h5>
                        <h1>{(data) ? replaceLast(data.cost.toFixed(2), '.', ',') : "calculando..."} ₲</h1>
                        <h1>{(data) ? data.consumption : 'calculando...'} wH</h1>
                    </div>
                    <div className="col">
                        <h5>Ayer</h5>
                        <h1>0 ₲</h1>
                        <h1>0 W</h1>
                    </div>
                    <div className="col">
                        <h5>Este mes</h5>
                        <h1>{(data) ? replaceLast(data.cost.toFixed(2), '.', ',') : "calculando..."} ₲</h1>
                        <h1>{(data) ? (data.consumption / 1000).toFixed(0) : 'calculando...'} KwH</h1>
                    </div>
                    <div className="col">
                        <h5>El mes pasado</h5>
                        <h1>0 ₲</h1>
                        <h1>0 kW</h1>
                    </div>
                    <div className="col">
                        <h5>Tarifa domicilio particular por KwH</h5>
                        <h1>402,000 ₲</h1>
                    </div>
                    <div className="col">
                        <h5>Tarifa industrial por KwH</h5>
                        <h1>334,798 ₲</h1>
                    </div>
                </div>
            </div>

            <h1 className="chart-title">Consumo de las ultimas 24hs</h1>
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
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="name" />
                    <YAxis />
                    <Tooltip />
                    <Line type="monotone" dataKey="uv" stroke="#8884d8" fill="#8884d8" />
                </LineChart>
            </ResponsiveContainer>
        </div>
    )
}

export default Dashboard;