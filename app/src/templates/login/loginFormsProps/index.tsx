import React from 'react';
import { useRouter } from 'next/router';

import { LoginFormsProps } from 'components/login';
import { usePageLoadingState } from 'hooks/usePageLoadingState';
import { useBasicValidationState, useVerificationCodeValidationState } from './useValidationState';
import { useVerificationCodeState } from 'hooks/useVerificationCodeState';
import { useLoginFormState } from './useLoingFormState';

import { EmailValidation, PasswordValidation, VerificationCodeValidataion } from 'domain/validations';
import { postLogin, postVerificateCode } from 'domain/apis/signin';
import { MyAccountContext } from 'hooks';

const emailValidator = new EmailValidation();
const passwordValidator = new PasswordValidation();
const verificataionCodeValidator = new VerificationCodeValidataion();

export function useLoginFormsProps() {
  const router = useRouter();
  const ac = React.useContext(MyAccountContext);
  React.useEffect(() => {
    if (ac.myAccount.has) {
      router.push('/');
    }
  }, [ac, router]);

  const pageState = usePageLoadingState<LoginFormsProps['switchContext']['content']>('basic', false);
  const basicState = useLoginFormState();
  const codeState = useVerificationCodeState();
  const codeValidationState = useVerificationCodeValidationState();
  const basicValidationState = useBasicValidationState();

  const props: LoginFormsProps = {
    switchContext: {
      content: pageState.currentPage,
    },
    basic: {
      email: {
        input: {
          value: basicState.email,
          onChange: (e) => basicState.setEmail(e.target.value),
          readOnly: pageState.isLoading,
        },
        description: {
          value: emailValidator.InputSuggestion,
        },
        coution: {
          value: basicValidationState.validationState.password,
        },
      },
      password: {
        input: {
          value: basicState.password,
          onChange: (e) => basicState.setPassword(e.target.value),
          readOnly: pageState.isLoading,
        },
        description: {
          value: passwordValidator.InputSuggestion,
        },
        coution: {
          value: basicValidationState.validationState.email,
        },
      },
      sendButton: {
        active: !pageState.isLoading,
        onClick: () => {
          // expecte code page after user-account is created.
          const page = pageState.presetPageAndStartLoading('code');
          const validatedEmail = emailValidator.validate(basicState.email);
          const validatedPassword = passwordValidator.validate(basicState.password);
          if (!validatedEmail.isValid || !validatedPassword.isValid) {
            const split = '';
            basicValidationState.setBasicValidationState(
              validatedEmail.getError(split),
              validatedPassword.getError(split),
            );
            page.resetCurrentPage();
            return;
          }
          // checked from input.

          postLogin(basicState.email, basicState.password)
            .then((userId) => {
              if (!userId) throw new Error();
              codeState.setUserId(userId);
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

          postVerificateCode(codeState.userId, codeState.code)
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
