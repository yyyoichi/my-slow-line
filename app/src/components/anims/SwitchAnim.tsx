import React from 'react';
import { FadeAnim, FadeAnimProps } from './FadeAnim';
import { DivProps } from 'components/atoms/type';

const SwitchContentContext = React.createContext<unknown>('');

// the context wrap contexts

export type SwitchAnimContextProps<T extends string> = Pick<DivProps, 'children'> & { content: T };

export const SwitchAnimContext = <T extends string>(props: SwitchAnimContextProps<T>) => {
  const [state, setState] = React.useState<T>(props.content);

  React.useEffect(() => setState(props.content), [props.content]);

  return <SwitchContentContext.Provider value={state}>{props.children}</SwitchContentContext.Provider>;
};

// contents

export type SwitchAnimContentProps<T extends string> = Pick<FadeAnimProps, 'children'> & {
  content: T;
};

export const SwitchAnimContent = <T extends string>({ content, children }: SwitchAnimContentProps<T>) => {
  const cxt = React.useContext(SwitchContentContext);
  const eq = cxt === content;
  const fadeProps: FadeAnimProps = {
    in: eq,
    children,
  };
  return <FadeAnim {...fadeProps} />;
};

export function Switch<T extends string>({ content }: { content: T }) {
  return {
    Content: ({ content: contentName, children }: SwitchAnimContextProps<T>) => {
      const eq = contentName === content;
      const fadeProps: FadeAnimProps = {
        in: eq,
        children,
      };
      return <FadeAnim {...fadeProps} />;
    },
  };
}
