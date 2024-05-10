import {
  Body,
  Button,
  Container,
  Column,
  Head,
  Heading,
  Hr,
  Html,
  Img,
  Link,
  Preview,
  Row,
  Section,
  Text,
} from "@react-email/components";
import { Tailwind } from "@react-email/tailwind";
import * as React from "react";
export const VercelInviteUserEmail = () => {
  return (
    <Html>
      <Head />
      <Tailwind>
        <Body className="bg-white my-auto mx-auto text-center font-sans px-2">
          <Container className="border border-solid border-[#eaeaea] rounded my-[40px] mx-auto p-[20px] max-w-[465px]">
            <Text className="text-black text-[20px] font-bold leading-[24px]">
              Selam var_name!,
            </Text>
            <Text className="text-black text-[14px] leading-[24px]">
              Ödevinin gitmesi için yollayanın sen olduğunu doğrulamalıyız!
            </Text>
            <Text className="text-black text-[20px] font-bold leading-[24px]">Doğrulama Kodun:</Text>
            <Text className="w-full text-center font-bold text-[20px] leading-[24px] text-blue-500 tracking-[0.5em] uppercase">var_code</Text>
            <Hr className="border border-solid border-[#eaeaea] my-[26px] mx-0 w-full" />
            <Link href="https://bento.me/haume" className="text-blue-500 text-[14px] leading-[24px]">
              by Haume
            </Link>
          </Container>
        </Body>
      </Tailwind>
    </Html>
  );
};
export default VercelInviteUserEmail;
