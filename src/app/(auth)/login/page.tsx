import { AuthCardHeader, AuthFooterForm } from "@/components/auth/card_content";
import LoginForm from "@/components/auth/login/login_form";

export default function Page() {
  return (
    <>
      <AuthCardHeader
        description="Enter your email and password"
        title="Welcome back!!"
      />
      <LoginForm />
      <AuthFooterForm
        href="/register"
        title="Don`t have an account?"
        title_href="Sign up"
      />
    </>
  );
}
