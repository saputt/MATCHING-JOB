import { IsEmail, IsString } from "class-validator"

export class RegisterDto {
    @IsString()
    username : string

    @IsString()
    @IsEmail()
    email : string

    @IsString()
    password : string
}