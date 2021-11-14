import { Navigate } from 'react-router-dom'
import { useState, useEffect } from "react";

import { Spin } from 'antd';

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

const loadDashboard = async (jwt, uci) => {
    if (!jwt || !uci) return { error: 401 }

    /**
    * Parametros del request
    * **/
    const params = {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${jwt}`,
            'mode': 'cors'
        },
    };

    /**
     * Llamada a nuestro servidor para que nos devuelve el resumen de uso
     * correspondiente a el usuario con C.I
     * **/
    const response = await fetch(`http://192.168.0.12:8080/usage/resume/${uci}`, params);

    if (!response.ok) return { error: response.status }

    return { data: "" }
}

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
            if (res.error) setError(res.error)

            setLoading(false);
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

    return (
        <h1>Dashboard</h1>
    )
}

export default Dashboard;