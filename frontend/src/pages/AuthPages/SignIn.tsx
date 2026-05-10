import PageMeta from "../../components/common/PageMeta";
import AuthLayout from "./AuthPageLayout";
import SignInForm from "../../components/auth/SignInForm";

export default function SignIn() {
  return (
    <>
      <PageMeta
        title="Sign In | my Chatbot Life"
        description="This is the Sign In page for my Chatbot Life"
      />
      <AuthLayout>
        <SignInForm />
      </AuthLayout>
    </>
  );
}
