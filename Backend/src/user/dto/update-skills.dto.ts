import { IsArray, IsNotEmpty, IsString } from "class-validator";

export class UpdateSkillsDto {
    @IsArray()
    @IsNotEmpty()
    @IsString({each : true})
    skills : string[]
}