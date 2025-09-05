import { type ReactNode } from "react";

interface CenteredLayoutProps {
  children: ReactNode;
}
const CenteredLayout = ({ children }: CenteredLayoutProps) => {
  return <div className="flex w-screen h-screen items-center place-content-center">{children}</div>;
};

export default CenteredLayout;
