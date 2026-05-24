import { ArgumentMetadata, Injectable, NotFoundException, PipeTransform } from "@nestjs/common";
import { UserRepository } from "../user.repository";
import { UserProfileResponse } from "../interfaces/user-response.interface";

@Injectable()
export class UserExist implements PipeTransform {
    constructor(private userRepo : UserRepository) {}

    async transform(userId : string) {
        const isUserExist = await this.userRepo.FindUserById(userId)

        if (!isUserExist) {
            throw new NotFoundException("User not exist")
        }

        return isUserExist
    }
}