import React from 'react';
import { ToggleViewIcon } from 'components/atoms/icons';
import { ToggleViewIconProps } from 'components/atoms/icons/ToggleView';
import {
  UiInputProps,
  UiInputDescriptionProps,
  UiInputCautionProps,
  UiInputLabelProps,
  UiInputLabel,
  UiInputDescription,
  UiInputCaution,
  UiInput,
} from 'components/atoms/input';
import { useState } from 'react';
import { NonNullablePick } from 'components';

export type PasswordFieldProps = {
  input: NonNullablePick<UiInputProps, 'value' | 'onChange' | 'readOnly'>;
};

export const PasswordField = (props: PasswordFieldProps) => {
  const [visiblePassword, setVisiblePassword] = useState<boolean>(false);
  const passwordToggleViewProps: PasswordToggleViewProps = {
    icon: {
      view: !visiblePassword,
      onClick: () => setVisiblePassword((v) => !v),
    },
  };
  const inputProps: UiInputProps = {
    type: visiblePassword ? 'password' : 'text',
    className: visiblePassword ? 'text-[1.3rem]' : '',
    ...props.input,
  };
  return (
    <div className='flex items-center bg-my-white pr-2'>
      <UiInput {...inputProps} />
      <PasswordToggleView {...passwordToggleViewProps} />
    </div>
  );
};

type PasswordToggleViewProps = {
  icon: Pick<ToggleViewIconProps, 'view' | 'onClick'>;
};
const PasswordToggleView = (props: PasswordToggleViewProps) => {
  const toggleViewIconProps: ToggleViewIconProps = {
    width: 22,
    height: 22,
    ...props.icon,
  };
  return <ToggleViewIcon {...toggleViewIconProps} />;
};

export type PasswordFrameProps = {
  description: Pick<UiInputDescriptionProps, 'value'>;
  coution: Pick<UiInputCautionProps, 'value'>;
} & PasswordFieldProps;

// password components
export const PasswordFrame = (props: PasswordFrameProps) => {
  const labelProps: UiInputLabelProps = {
    value: 'Password',
  };
  const descriptionProps: UiInputDescriptionProps = {
    ...props.description,
  };
  const coutionProps: UiInputDescriptionProps = {
    ...props.coution,
  };

  return (
    <>
      <UiInputLabel {...labelProps} />
      <UiInputDescription {...descriptionProps} />
      <UiInputCaution {...coutionProps} />
      <PasswordField {...props} />
    </>
  );
};
