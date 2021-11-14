import './App.css';
import logo from './logo.svg';

import React from 'react';

import { Layout, Divider, Button } from 'antd';
import { BrowserRouter as Router, Route, Routes, useLocation } from "react-router-dom";

import Landing from './Landing';
import Dashboard from './dashboard/Dashboard';
import Login from './login/Login';


/**
 * Estilo para posicionar elementos adentro de manera horizontal
 * **/
const rowStyle = {
  display: 'flex',
  flexDirection: 'row', // Posisionar elementos horizontalmente
  alignItems: 'center', // Centrar verticalmente
};

/**
 * Logo de la App
 * **/
function AppLogo() {
  const logoStyle = {
    height: '40px',
    marginRight: '16px' // Dejar margen hacia el lado derecho
  };
  return (
    <a href="/" style={rowStyle}>
      <img style={logoStyle} src={logo} alt="" />
      <h2 style={{ margin: 0 }}>PowerMeter</h2>
    </a>
  );
}

/**
 * Encabezado de pagina
 * **/
function AppHeader() {
  const location = useLocation();
  return (
    <div className="app-header">
      <AppLogo />
      <div style={rowStyle}>
        {
          (!location.pathname.endsWith("/login")) ?
            <Button id="btn" type="primary" href="/login">
              Iniciar Sesion
            </Button>
            : <></>
        }
      </div>
    </div>
  );
}

/**
* Cuerpo de la pagina
* **/
function AppContent() {
  return (
    <div className="app-content-body">
      <Routes>
        <Route path="" element={<Landing />} />
        <Route path="login" element={<Login />} />
        <Route path="dashboard" element={<Dashboard />} />
      </Routes>
    </div>
  );
}

/**
* Parte de abajo de la pagina
* **/
function AppFooter() {
  return (
    <div className="app-footer">
      <Divider />
      <p>
        Proyecto realizado con fines de aprendizaje para la <a href="https://www.unida.edu.py/facultades/facultad-de-ingenieria/ingenieria-informatica/">Facultad de Ing. Informatica UNIDA</a>
      </p>
      <p>Colaboradores:</p>
      <p>
        <a href="mailto:juanhr454@gmail.com">Juan Herrera</a> • <a href="mailto:deboraareliescobar@gmail.com">Debora Escobar</a> • <a href="mailto:deliaortizservin@gmail.com">Delia Ortiz</a>
      </p>
      <p>
        <a href="https://github.com/imherrera/medidor-electronico">Ver codigo fuente en Github</a>
      </p>
    </div>
  );
}

/**
 * Contenedor de la pagina
 * **/
function App() {
  return (
    <Router>
      <Layout style={{ minHeight: '100vh' }}>
        <AppHeader />
        <AppContent />
        <AppFooter />
      </Layout>
    </Router>
  );
}

export default App;