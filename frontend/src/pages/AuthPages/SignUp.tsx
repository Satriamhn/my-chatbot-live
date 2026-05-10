import PageMeta from "../../components/common/PageMeta";
import AuthLayout from "./AuthPageLayout";
import SignUpForm from "../../components/auth/SignUpForm";

export default function SignUp() {
  return (
    <>
      <PageMeta
        title="Sign Up | my Chatbot Life"
        description="This is the Sign Up page for my Chatbot Life"
      />
      <AuthLayout>
        <SignUpForm />
      </AuthLayout>
    </>
  );
}
