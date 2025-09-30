import { AuthCardHeader, AuthFooterForm } from "@/components/auth/card_content";
import RegisterForm from "@/components/auth/register/register_form";

export default function Page() {
  return (
    <>
      <AuthCardHeader description="Create your account" title="Welcome!!" />
      <RegisterForm />
      <AuthFooterForm
        href="/login"
        title="Have an account?"
        title_href="Login"
      />
    </>
  );
}
