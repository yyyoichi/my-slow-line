import React from 'react';
import { SigninFormsProps } from 'components/signin';
import {
  EmailValidation,
  PasswordValidation,
  ConfirmPasswordValidataion,
  NameValidation,
  VerificationCodeValidataion,
} from 'domain/validations';
import { usePageLoadingState } from 'hooks/usePageLoadingState';
import { useSigninFormState } from './useSigninFormState';
import { useVerificationCodeState } from 'hooks/useVerificationCodeState';
import { useBasicValidationState, useVerificationCodeValidationState } from './useValidationState';
import { postSignin, postVerificateCode } from 'domain/apis';
import { useRouter } from 'next/router';
import { MyAccountContext } from 'hooks';

const emailValidator = new EmailValidation();
const passwordValidator = new PasswordValidation();
const confirmPassValidator = new ConfirmPasswordValidataion();
const nameValidator = new NameValidation();
const verificataionCodeValidator = new VerificationCodeValidataion();

export function useSigninFormsProps() {
  const ac = React.useContext(MyAccountContext);
  const router = useRouter();
  React.useEffect(() => {
    if (ac.myAccount.has) {
      router.push('/');
    }
  }, [ac, router]);

  const pageState = usePageLoadingState<SigninFormsProps['switchContext']['content']>('name', false);
  const basicState = useSigninFormState();
  const codeState = useVerificationCodeState();
  const codeValidationState = useVerificationCodeValidationState();
  const basicValidationState = useBasicValidationState();

  const props: SigninFormsProps = {
    switchContext: {
      content: pageState.currentPage,
    },
    userName: {
      name: {
        input: {
          value: basicState.name,
          onChange: (e) => basicState.setName(e.target.value),
          readOnly: pageState.isLoading,
        },
        description: {
          value: nameValidator.InputSuggestion,
        },
        coution: {
          value: basicValidationState.validationState.name,
        },
      },
      sendButton: {
        active: !pageState.isLoading,
        onClick: () => {
          const page = pageState.presetPageAndStartLoading('basic');

          const validatedName = nameValidator.validate(basicState.name);
          if (!validatedName.isValid) {
            basicValidationState.setNameState(validatedName.getError(' '));
            page.resetCurrentPage();
            return;
          }

          // checked form input.

          page.goToNextPage();
        },
      },
    },
    basic: {
      email: {
        input: {
          value: basicState.email,
          onChange: (e) => {
            basicState.setEmail(e.target.value);
          },
          readOnly: pageState.isLoading,
        },
        description: {
          value: emailValidator.InputSuggestion,
        },
        coution: {
          value: basicValidationState.validationState.email,
        },
      },
      password: {
        input: {
          value: basicState.password,
          onChange: (e) => {
            basicState.setPassword(e.target.value);
          },
          readOnly: pageState.isLoading,
        },
        description: {
          value: passwordValidator.InputSuggestion,
        },
        coution: {
          value: basicValidationState.validationState.password,
        },
      },
      confirmPassword: {
        input: {
          value: basicState.confirmPassword,
          onChange: (e) => {
            basicState.setConfirmPassword(e.target.value);
          },
          readOnly: pageState.isLoading,
        },
      },
      sendButton: {
        active: !pageState.isLoading,
        onClick: () => {
          // expecte code page after user-account is created.
          const page = pageState.presetPageAndStartLoading('code');

          const { email, password, confirmPassword } = basicState;
          const validatedEmail = emailValidator.validate(email);
          const validatedPass = passwordValidator.validate(password);
          const validatedConfPass = confirmPassValidator.validate([password, confirmPassword]);
          if (!validatedEmail.isValid || !validatedPass.isValid || !validatedConfPass.isValid) {
            const split = ' ';
            basicValidationState.setBasicValidationState(
              validatedEmail.getError(split),
              validatedPass.getError(split),
              validatedConfPass.getError(split),
            );
            page.resetCurrentPage();
            return;
          }
          // checked form input.

          postSignin(email, password, basicState.name)
            .then((jwt) => {
              if (!jwt) throw new Error();
              codeState.setJwt(jwt);
              page.goToNextPage();
            })
            .catch(() => {
              alert('Sorry, occurs unexpected errors. Please try agin.');
              page.resetCurrentPage();
            });
        },
      },
    },
    verificationCode: {
      code: {
        input: {
          value: codeState.code,
          onChange: (e) => codeState.setCode(e.target.value),
          readOnly: pageState.isLoading,
        },
        description: {
          value: verificataionCodeValidator.InputSuggestion,
        },
        coution: {
          value: codeValidationState.validationState,
        },
      },
      sendButton: {
        active: !pageState.isLoading,
        onClick: () => {
          const page = pageState.presetPageAndStartLoading('code');

          const validatedCode = verificataionCodeValidator.validate(codeState.code);
          if (!validatedCode.isValid) {
            const split = ' ';
            codeValidationState.setState(validatedCode.getError(split));
            page.resetCurrentPage();
            return;
          }

          // checked form input.

          postVerificateCode(codeState.jwt, codeState.code)
            .then(() => ac.pullMyAccount())
            .catch(() => {
              alert('Sorry, occurs unexpected errors. Please try agin.');
              page.resetCurrentPage();
            });
        },
      },
    },
  };
  return props;
}
