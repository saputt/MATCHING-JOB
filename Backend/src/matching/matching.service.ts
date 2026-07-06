import { BadRequestException, Injectable } from "@nestjs/common";
import { JobService } from "src/job/job.service";
import { UserService } from "src/user/user.service";

@Injectable()
export class MatchingService {
    constructor(
        private readonly userService : UserService,
        private readonly jobService : JobService
    ) {}

    async getTopRecomendation(userId : string) {
        await this.userService.isUserExist(userId)

        const embedding = await this.userService.getUserEmbedding(userId)

        if (!embedding) {
            throw new BadRequestException("User doesnt have embedding. setup preference first");
        }

        const topJobs = await this.jobService.findTopRecommendation(embedding)

        return topJobs
    }
}   