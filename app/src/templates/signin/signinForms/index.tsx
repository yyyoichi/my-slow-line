import React from 'react';
import { SigninForms, SigninFormsProps } from 'components/signin';

export default function Signin() {
  const [state, setState] = React.useState<SigninFormsProps['switchContext']['content']>('basic');
  const props: SigninFormsProps = {
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
      confirmPassword: {
        input: {
          value: 'confirm',
          onChange: () => console.log('r'),
          readOnly: true,
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
        onClick: () => setState('name'),
      },
    },
    userName: {
      name: {
        input: {
          value: 'name',
          onChange: () => console.log('click'),
          readOnly: false,
        },
        coution: {
          value: '',
        },
        description: {
          value: 'user account name',
        },
      },
      sendButton: {
        active: true,
        onClick: () => setState('basic'),
      },
    },
  };

  return <SigninForms {...props} />;
}
