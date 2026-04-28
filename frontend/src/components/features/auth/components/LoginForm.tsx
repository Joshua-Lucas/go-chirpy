import FormField from "../../../ui/FormField"
import { ArrowRightIcon } from "@heroicons/react/16/solid"
import PrimaryButton from "../../../ui/PrimaryButton"

function LoginForm() {
  return (
    <form className="space-y-lg">
      <FormField id="email" type="email" placeholder="name@company.com" label="Email Address" />
      <FormField id="password" type="password" placeholder="••••••••" label="Password" />
      <PrimaryButton type="submit" text="Login" Icon={ArrowRightIcon} />
    </form>

  )
}

export default LoginForm
