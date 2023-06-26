import React from 'react';
import { NonNullablePick } from 'components';
import FormFrame from 'components/frame/FormFrame';
import { LoadButton, LoadButtonProps } from 'components/frame/LoadButton';
import { SwitchAnimContent, SwitchAnimContext, SwitchAnimContextProps } from 'components/anims/SwitchAnim';

type AnimType = 'login' | 'pwa' | 'notification' | '';

export type HomePageProps = {
  switchContext: SwitchAnimContextProps<AnimType>;
  notification: ReqNotifierProps;
};
export const HomePage = ({ switchContext, ...props }: HomePageProps) => {
  const switchAnimContextProps: SwitchAnimContextProps<AnimType> = {
    ...switchContext,
  };
  return (
    <SwitchAnimContext {...switchAnimContextProps}>
      <SwitchAnimContent content='login'>
        <ReqLogin />
      </SwitchAnimContent>
      <SwitchAnimContent content='pwa'>
        <ReqPwa />
      </SwitchAnimContent>
      <SwitchAnimContent content='notification'>
        <ReqNotifier {...props.notification} />
      </SwitchAnimContent>
    </SwitchAnimContext>
  );
};

const ReqLogin = () => {
  return (
    <FormFrame.Container>
      <FormFrame.Content>Please Login!</FormFrame.Content>
    </FormFrame.Container>
  );
};

const ReqPwa = () => {
  return (
    <FormFrame.Container>
      <FormFrame.Content>
        <p>Please Add Home.</p>
        <p>This app must be added to the home screen.</p>
      </FormFrame.Content>
    </FormFrame.Container>
  );
};

type ReqNotifierProps = {
  notificationButton: NonNullablePick<LoadButtonProps, 'onClick' | 'active'>;
  tmpResult: string;
};
const ReqNotifier = (props: ReqNotifierProps) => {
  const buttonProps: LoadButtonProps = {
    color: 'yellow',
    children: 'Subscribe notifications',
    ...props.notificationButton,
  };
  return (
    <FormFrame.Container>
      <FormFrame.Content>
        <p>For secure communication, This app must allow notifications.</p>
        <LoadButton {...buttonProps} />
        <p>{props.tmpResult}</p>
      </FormFrame.Content>
    </FormFrame.Container>
  );
};
