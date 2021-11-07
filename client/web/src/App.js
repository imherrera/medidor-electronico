import { useState } from 'react';
import { Layout, Menu, Button } from 'antd';
import { BrowserRouter as Router, Switch, Route, useLocation } from "react-router-dom";
const { Header, Content, Footer } = Layout;

function Login() {
  const error = null;

  function submit(credentials) {
    // En esta funcion enviamos los credenciales (Correo, Contraseña) al servidor
  }

  return <h1>Pantalla de inicio de sesion</h1>
}

function Dashboard() {

  return <h1>Pantalla de estadisticas</h1>
}

const routes = [
  {
    title: 'Inicio',
    path: '/',
    main: () => <Dashboard />
  },
  {
    title: 'Iniciar sesion',
    path: '/login',
    main: () => <Login />
  },
  {
    title: '404',
    path: '/*',
    main: () => <h1>Pagina no encontrada</h1>
  }
];

function App() {
  return (
    <Router>
      <Layout style={{ minHeight: '100vh' }}>
        <Header>
          <Menu theme="dark" mode="horizontal">          
            <Menu.Item>
              Iniciar Sesion
            </Menu.Item>
          </Menu>
        </Header>
      </Layout>
    </Router>
  );
}

export default App;
