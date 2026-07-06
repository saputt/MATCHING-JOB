import { createParamDecorator, ExecutionContext } from "@nestjs/common";

interface UserPayload {
    id : string;
    email : string
}

export const GetUser = createParamDecorator(
    (data : keyof UserPayload | undefined, ctx: ExecutionContext) => {
        const req = ctx.switchToHttp().getRequest()
        const user = req.user as UserPayload
        return data ? user[data] : user
    }
)