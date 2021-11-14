import './Login.css';
import logo from '../logo.svg';

import { useState } from 'react';

import { Button, Input, Form, Checkbox, Alert } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';

const postLoginForm = async (form) => {
    /**
     * Credenciales de logeo del usuario
     * **/
    const credentials = {
        uci: form['user_ci'],
        pass: form['password']
    };

    /**
     * Parametros del request
     * **/
    const params = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(credentials) // serializamos el obj a un string json
    };

    /**
     * Llamada a nuestro servidor para logearnos
     * **/
    const response = await fetch('http://192.168.0.12:8080/login', params);

    /**
     * Checkeamos que el servidor acepto las credenciales como validas
     * **/
    if (!response.ok) return { error: response.status }

    return await response.json()
};

function LoginFormError() {
    return (
        <Alert
            style={{ marginBottom: "2em" }}
            className="form-error"
            message="Error de acceso"
            description="Tus credenciales no son correctos o validos"
            type="error"
            showIcon
            closable
        />
    );
}


function LoginForm() {
    const [error, setError] = useState(false);

    function onFinish(form) {
        postLoginForm(form).then(res => {
            if (res.error === 401) {
                setError(true);
                localStorage.setItem('power-meter-jwt', null);
            } else {
                setError(false);
                /**
                 * Guardamos el jwt en local storage
                 * 
                 * NOTA: esto no se debe hacer en una applicacion real
                 * ***/
                localStorage.setItem('power-meter-jwt', res.tkn);
                localStorage.setItem('power-meter-uci', form['user_ci']);
            }
        }).catch(err => {
            console.error("Error submiting login form:", err)
        })
    };

    return (
        <>
            {(error) ? <LoginFormError /> : <></>}
            <Form
                className="login-form"
                name="basic"
                wrapperCol={{ span: 100 }}
                initialValues={{ remember: true }}
                onFinish={onFinish}
                autoComplete="off"
            >
                <Form.Item
                    name="user_ci"
                    rules={[
                        {
                            required: true,
                            message: 'Debes ingresar tu C.I',
                        },
                    ]}
                >
                    <Input prefix={<UserOutlined />} placeholder="Numero de C.I" />
                </Form.Item>

                <Form.Item
                    name="password"
                    rules={[
                        {
                            required: true,
                            message: 'Debes ingresar tu contraseña',
                        },
                    ]}
                >
                    <Input.Password prefix={<LockOutlined />} placeholder="Contraseña" />
                </Form.Item>

                <Form.Item
                    name="remember"
                    valuePropName="checked"
                    wrapperCol={{
                        offset: 100,
                        span: 16,
                    }}
                >
                    <Checkbox>Recordarme</Checkbox>
                </Form.Item>

                <Form.Item style={{ alignSelf: 'center' }}>
                    <Button id="btn" type="primary" htmlType="submit">
                        Continuar
                    </Button>
                </Form.Item>
            </Form>
        </>
    );
};


function Login() {
    return (
        <div className="login-layout">
            <div className="login-container">
                <img className="logo" src={logo} alt="" />
                <LoginForm />
            </div>
        </div>
    )
}

export default Login;