import { IsArray, IsNotEmpty, IsOptional, IsString } from "class-validator";

export class UpdatePreferencesDto {
    @IsArray()
    @IsNotEmpty()
    @IsString({each : true})
    skills : string[]

    @IsString()
    @IsOptional()
    titleInterest : string
}