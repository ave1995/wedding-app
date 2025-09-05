import { type ReactNode } from "react";

interface CenteredLayoutProps {
  children: ReactNode;
}
const CenteredLayout = ({ children }: CenteredLayoutProps) => {
  return <div className="flex w-screen h-svh justify-center items-center">{children}</div>;
};

export default CenteredLayout;
