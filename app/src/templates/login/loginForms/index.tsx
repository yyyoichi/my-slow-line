import React from 'react';
import { LoginForms, LoginFormsProps } from 'components/login';

export default function Login() {
  const [state, setState] = React.useState<LoginFormsProps['switchContext']['content']>('basic');
  const props: LoginFormsProps = {
    switchContext: {
      content: state,
    },
    basic: {
      password: {
        input: {
          value: 'text',
          onChange: () => setState('code'),
          readOnly: false,
        },
        description: {
          value: 'prese entry',
        },
        coution: {
          value: 'nil',
        },
      },
      email: {
        input: {
          value: 'text',
          onChange: (e) => {
            console.log(e.target.value);
          },
          readOnly: true,
        },
        description: {
          value: 'prese entry',
        },
        coution: {
          value: 'nil',
        },
      },
      sendButton: {
        active: true,
        onClick: () => setState('code'),
      },
    },
    verificationCode: {
      code: {
        input: {
          value: '000000',
          onChange: () => console.log('click'),
          readOnly: false,
        },
        description: {
          value: '6 code',
        },
        coution: {
          value: '',
        },
      },
      sendButton: {
        active: true,
        onClick: () => setState('basic'),
      },
    },
  };

  return <LoginForms {...props} />;
}
