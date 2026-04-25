import BackgroundGradient from "../components/BackgroundGradiant"
import AuthHeader from "../components/features/auth/components/AuthHeader"
import AuthLayout from "../components/layout/authlayout"

function LoginPage() {

  return (
    <>
      <BackgroundGradient />
      <AuthLayout>
        <AuthHeader />
        <p>Hello World</p>
      </AuthLayout>
    </>
  )

}

export default LoginPage
