import './Dashboard.css';

import { Navigate } from 'react-router-dom'
import { useState, useEffect } from "react";

import { Spin } from 'antd';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { isMobile } from '../utils';




function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

const kwhHomePrice = 402.000;
const kwhBussinessPrice = 334.798

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

/*const data = [
    {
        name: 'Page A',
        uv: 4000,
        pv: 2400,
        amt: 2400,
    },
    {
        name: 'Page B',
        uv: 3000,
        pv: 1398,
        amt: 2210,
    },
    {
        name: 'Page C',
        uv: 2000,
        pv: 9800,
        amt: 2290,
    },
    {
        name: 'Page D',
        uv: 2780,
        pv: 3908,
        amt: 2000,
    },
    {
        name: 'Page E',
        uv: 1890,
        pv: 4800,
        amt: 2181,
    },
    {
        name: 'Page F',
        uv: 2390,
        pv: 3800,
        amt: 2500,
    },
    {
        name: new Date().toLocaleDateString(),
        uv: 3490,
        pv: 4300,
        amt: 2100,
    },
];*/

const padding = isMobile() ? 16 : 66

function Dashboard() {
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [data, setData] = useState(null);

    useEffect(() => {
        setLoading(true);
        /**
         * Leemos el token guardado, sera nulo si el usuario no tiene una sesion abierta
         * **/
        const token = localStorage.getItem('power-meter-jwt');
        const uci = localStorage.getItem('power-meter-uci');

        loadDashboard(token, uci).then(res => {
            setLoading(false);
            if (res.error) {
                setError(res.error)
                return
            }

            const graphData = [];
            res.forEach(e => {
                const cost = (1000 / e.watt_hour) * kwhHomePrice;
                const name = "KwH"
                graphData.push({
                    name: 'W/H',
                    pv: new Date(e.date).toTimeString(),
                    uv: e.watt_hour,
                    cost: cost.toFixed(2),
                })
            });
            setData(graphData);

            console.log("Response from api: ", res)
        }).catch(err => {
            console.error("Failure on getting resume:", err);
        })

    }, []);

    if (loading) {
        return (
            <div style={{
                width: '100%',
                height: '70vh',
                display: 'flex',
                flexDirection: 'column',
                justifyContent: 'center'
            }}>
                <Spin style={{ alignSelf: 'center' }} tip="Cargando..." size="large"></Spin>
            </div>
        )
    }

    /**
     * Redirigimos al usuario a la pantalla de inicio de sesion
     * **/
    if (error === 401) {
        return <Navigate to='/login' />
    }

    console.log("data:", data)
    return (
        <div className="dashboard-layout">
            <h1>Consumo Diario</h1>
            <LineChart
                className="chart"
                data={data}
                width={((isMobile()) ? 1 : 0.80) * window.outerWidth}
                height={0.60 * window.innerHeight}
                margin={{
                    top: padding,
                    right: padding,
                    left: padding / 2,
                    bottom: padding,
                }}
            >
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis dataKey="cost" />
                <Tooltip />
                <Legend />
                <Line type="monotone" dataKey="pv" stroke="#8884d8" activeDot={{ r: 8 }} />
                <Line type="monotone" dataKey="uv" stroke="#82ca9d" />
            </LineChart>
        </div>
    )
}

export default Dashboard;