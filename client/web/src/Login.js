import './Login.css';
import logo from './logo.svg';

import { Button, Input, Space, Form, Checkbox } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';

const Demo = () => {
    const onFinish = (values) => {
        console.log('Success:', values);
    };

    const onFinishFailed = (errorInfo) => {
        console.log('Failed:', errorInfo);
    };

    return (
        <Form
            className="login-form"
            name="basic"
            wrapperCol={{ span: 100 }}
            initialValues={{ remember: true }}
            onFinish={onFinish}
            onFinishFailed={onFinishFailed}
            autoComplete="off"
        >
            <Form.Item
                name="username"
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
    );
};

function LoginHeader() {
    return (
        <img className="logo" src={logo} alt="" />
    );
}

function LoginForm() {
    return (
        <Space direction="vertical" size="small">
            <Input prefix={<UserOutlined />} placeholder="Numero de C.I" style={{ width: "50em" }} />
            <Input.Password prefix={<LockOutlined />} placeholder="Contraseña" />
        </Space>
    )
}

function LoginContainer(props) {
    return (
        <div className="login-container">
            {props.children}
        </div>
    );
}

function Login() {
    return (
        <div className="login-layout">
            <LoginContainer>
                <LoginHeader />
                <Demo />
            </LoginContainer>
        </div>
    )
}



export default Login;